package main

import (
	"context"
	"os"

	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/app/server"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/config"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	cfg, err := config.New(env)
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.AuthenticationService.LogLevel)

	if err := server.Init(context.Background(), cfg, logger); err != nil {
		logger.Fatal("server initialization failed", log.Error(err))
	}
}
