package applicator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/repo"
	"image-gallery/internal/gallery/service"
	"image-gallery/internal/gallery/transport"
	"image-gallery/internal/gallery/worker"

	"image-gallery/internal/gallery/config"
	"image-gallery/internal/gallery/db"
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

	workerCount := 3
	tasks := make(chan entity.Image, 1)
	result := make(chan string, 10)
	out := make(<-chan string, 10)

	repository := repo.NewRepository(database.GetDB())
	Worker := worker.NewWorker(workerCount, tasks, result, out, *repository)
	authGrpcTransport := transport.NewAuthGrpcTransport(cfg.Transport.AuthGrpc)
	userGrpcTransport := transport.NewUserGrpcTransport(cfg.Transport.UserGrpc)
	galleryService := service.NewService(*repository, log, authGrpcTransport, userGrpcTransport, Worker)
	galleryService.WorkerRunInService()
	handler := service.NewHandler(galleryService, Worker)

	service.InitRouters(handler, r)

	serverPort := cfg.Server.Port
	addr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("Error in running server: %s", err)
	}
}
