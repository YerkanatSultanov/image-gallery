package main

import (
	"go.uber.org/zap"
	"image-gallery/internal/auth/applicator"
	"image-gallery/internal/auth/config"
)

func main() {

	logger, _ := zap.NewProduction()
	//nolint:all
	defer logger.Sync()
	l := logger.Sugar()
	l = l.With(zap.String("app", "auth-service"))
	conf, err := config.LoadConfig("config/auth")
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}
	app := applicator.NewApplicator(l, conf)
	app.Run()
}
