package applicator

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/auth/config"
	"image-gallery/internal/auth/db"
	"image-gallery/internal/auth/repo"
	"image-gallery/internal/auth/service"
	"image-gallery/internal/auth/service/grpc"
	"image-gallery/internal/auth/transport"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Applicator struct {
	logger *zap.SugaredLogger
	config config.Config
	server *http.Server
	wg     sync.WaitGroup
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

	a.server = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Fatalf("Error in running server: %s", err)
		}
	}()

	a.gracefulShutdown()

	a.wg.Wait()
	a.logger.Info("Server gracefully stopped")
}

func (a *Applicator) gracefulShutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	<-signalCh

	a.logger.Info("Received shutdown signal. Shutting down...")

	err := a.server.Shutdown(context.Background())
	if err != nil {
		a.logger.Errorf("Error during server shutdown: %s", err)
	}
}
