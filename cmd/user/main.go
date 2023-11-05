package main

import (
	"fmt"
	"github.com/spf13/viper"
	"image-gallery/internal/user/config"
	"image-gallery/internal/user/db"
	"image-gallery/internal/user/repo"
	"image-gallery/internal/user/service"
	"log"
)

func main() {
	cfg, err := loadConfig("config/user")

	if err != nil {
		return
	}

	dbConn, err := db.NewDataBase(cfg.Database)

	if err != nil {
		log.Fatalf("could not initialize database connection %s", err)
	}

	userRep := repo.NewRepository(dbConn.GetDB())
	userService := service.NewService(userRep)
	userHandler := service.NewHandler(userService)

	service.InitRouters(userHandler)

	serverPort := cfg.Server.Port
	addr := fmt.Sprintf("0.0.0.0:%d", serverPort)

	err = service.Start(addr)
	if err != nil {
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
