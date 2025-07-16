package main

import (
	"GIN/configs"
	"GIN/internal/app"
	"log"
)

func main() {
	cfg := configs.LoadConfig()
	if cfg == nil {
		log.Fatalf("failed to load config")
	}

	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}