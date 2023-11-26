package main

import (
	"car-pooling-challenge/internal/app"
	"car-pooling-challenge/pkg/config"
	"car-pooling-challenge/pkg/logger"
)

func main() {

	// Configuration
	cfg, err := config.Load(".env.conf")

	if err != nil {
		logger.NewLogger("Main").Fatalf("Config error: %s", err)
	}

	// Run
	app := app.New(cfg)

	app.Run()
}
