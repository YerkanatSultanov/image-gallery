package applicator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/auth/config"
	"image-gallery/internal/auth/db"
	"image-gallery/internal/auth/repo"
	"image-gallery/internal/auth/service"
	"image-gallery/internal/auth/service/grpc"
	"image-gallery/internal/auth/transport"
)

type Applicator struct {
	logger *zap.SugaredLogger
	config config.Config
}

func NewApplicator(logger *zap.SugaredLogger, cfg config.Config) *Applicator {
	return &Applicator{
		logger: logger,
		config: cfg,
	}
}

func (a *Applicator) Run() {
	r := gin.Default()
	cfg := a.config
	log := a.logger
	database, err := db.NewDataBase(cfg.Database)
	if err != nil {
		log.Fatalf("Error to connect database: %s", err)
	}

	userRep := repo.NewRepository(database.GetDB())
	userGrpcTransport := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)
	trans := transport.NewTransport(cfg.Transport.User, log)
	userService := grpc.NewService(userRep, trans, userGrpcTransport, log, cfg)
	userHandler := service.NewHandler(*userService)

	grpcServer := grpc.NewServer(cfg.GrpcServer.Port, userService)
	err = grpcServer.Start()
	if err != nil {
		log.Panicf("failed to start grpc-server err: %v", err)
	}

	service.InitRouters(userHandler, r)

	serverPort := cfg.Server.Port
	addr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("Error in running server: %s", err)
	}
}
