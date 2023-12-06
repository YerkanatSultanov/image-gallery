package main

import (
	"go.uber.org/zap"
	"image-gallery/internal/gallery/applicator"
	"image-gallery/internal/gallery/config"
)

// @title						Image-Gallery service
// @version					1.0
// @description				Image-Gallery service
// @termsOfService				http://swagger.io/terms/
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8080
// @securityDefinitions.apikey	BearerAuth
// @type						apiKey
// @name						Authorization
// @in							header
// @schemes					http
func main() {
	logger, _ := zap.NewProduction()
	//nolint:all
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
