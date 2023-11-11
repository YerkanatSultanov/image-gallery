package applicator

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image-gallery/internal/user/config"
	"image-gallery/internal/user/db"
	"image-gallery/internal/user/repo"
	"image-gallery/internal/user/service"
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
	userService := service.NewService(userRep)
	userHandler := service.NewHandler(userService)

	service.InitRouters(userHandler, r)

	serverPort := cfg.Server.Port
	addr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	err = r.Run(addr)
	if err != nil {
		log.Fatalf("Error in running server: %s", err)
	}
}
