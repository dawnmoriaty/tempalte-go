package main

import (
	"GIN/configs"
	_ "GIN/docs"
	"GIN/internal/app"
	"GIN/pkg/logger"
	"go.uber.org/zap"
)

// go
// @title           E-commerce API
// @version         1.0
// @description     API cho hệ thống Auth (JWT + Redis) dùng Gin.
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Type "Bearer <token>" (without quotes). Example: Bearer eyJhbGciOi...
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
func main() {
	logger.InitLogger()
	zap.S().Info("Starting application...")
	cfg := configs.LoadConfig()
	if cfg == nil {
		zap.S().Fatal("failed to load config")
	}

	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		zap.S().Fatalf("failed to run application: %v", err)
	}
}
