package applicator

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/gallery/entity"
	"image-gallery/internal/gallery/repo"
	"image-gallery/internal/gallery/service"
	"image-gallery/internal/gallery/transport"
	"image-gallery/internal/gallery/worker"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"image-gallery/internal/gallery/config"
	"image-gallery/internal/gallery/db"
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
