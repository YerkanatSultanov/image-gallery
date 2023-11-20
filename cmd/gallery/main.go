package main

import (
	"go.uber.org/zap"
	"image-gallery/internal/gallery/applicator"
	"image-gallery/internal/gallery/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	l := logger.Sugar()
	l = l.With(zap.String("app", "gallery-service"))
	conf, err := config.LoadConfig("config/gallery")
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}
	app := applicator.NewApplicator(l, conf)
	app.Run()
}
