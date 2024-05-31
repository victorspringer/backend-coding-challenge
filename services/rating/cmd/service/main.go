package main

import (
	"context"
	"os"

	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/app/server"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/log"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	cfg, err := config.New(env)
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.RatingService.LogLevel)

	if err := server.Init(context.Background(), cfg, logger); err != nil {
		logger.Fatal("server initialization failed", log.Error(err))
	}
}
