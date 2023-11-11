package main

import (
	"go.uber.org/zap"
	"image-gallery/internal/user/applicator"
	"image-gallery/internal/user/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	l := logger.Sugar()
	l = l.With(zap.String("app", "user-service"))
	conf, err := config.LoadConfig("config/user")
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}
	app := applicator.NewApplicator(l, conf)
	app.Run()
}
