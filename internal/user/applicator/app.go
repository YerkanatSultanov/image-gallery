package applicator

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/kafka"
	"image-gallery/internal/user/config"
	"image-gallery/internal/user/controller/consumer"
	"image-gallery/internal/user/db"
	"image-gallery/internal/user/repo"
	"image-gallery/internal/user/service"
	"image-gallery/internal/user/service/grpc"
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

	userVerificationProducer, err := kafka.NewProducer(cfg.Kafka)
	if err != nil {
		log.Panicf("failed NewProducer err: %v", err)
	}

	userVerificationConsumerCallback := consumer.NewUserVerificationCallback(log)

	userVerificationConsumer, err := kafka.NewConsumer(log, cfg.Kafka, userVerificationConsumerCallback)
	if err != nil {
		log.Panicf("failed NewConsumer err: %v", err)
	}

	go userVerificationConsumer.Start()

	userRep := repo.NewRepository(database.GetDB())
	userService := grpc.NewService(userRep, log, userVerificationProducer)
	userHandler := service.NewHandler(userService)

	grpcServer := grpc.NewServer(cfg.GrpcServer.Port, userService)
	err = grpcServer.Start()
	if err != nil {
		log.Panicf("failed to start grpc-server err: %v", err)
	}

	defer grpcServer.Close()

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
