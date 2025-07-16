package main

import (
	"GIN/configs"
	"GIN/internal/app"
	"GIN/pkg/logger"
	"go.uber.org/zap"
)

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