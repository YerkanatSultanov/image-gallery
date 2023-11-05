package main

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image-gallery/internal/auth/config"
	"image-gallery/internal/auth/db"
	"image-gallery/internal/auth/repo"
	"image-gallery/internal/auth/service"
	"image-gallery/internal/auth/transport"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	cfg, err := loadConfig("config/auth")

	if err != nil {
		logger.Info("Failed at reading config")
		return
	}

	dbConn, err := db.NewDataBase(cfg.Database)

	if err != nil {
		logger.Info("could not initialize database connection")
	}

	rep := repo.NewRepository(dbConn.GetDB())

	trans := transport.NewTransport(cfg.Transport.User, logger)
	userService := service.NewService(rep, trans, logger)
	userHandler := service.NewHandler(userService, rep)

	service.InitRouters(userHandler)

	serverPort := cfg.Server.Port
	addr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	err = service.Start(addr)
	if err != nil {
		logger.Info("can not connect to address")
		return
	}
}

func loadConfig(path string) (config config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")

	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to ReadInConfig err: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to Unmarshal config err: %w", err)
	}

	return config, nil
}
