package main

import (
	"context"
	"os"

	"github.com/victorspringer/backend-coding-challenge/services/user/internal/app/server"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	cfg, err := config.New(env)
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.UserService.LogLevel)

	if err := server.Init(context.Background(), cfg, logger); err != nil {
		logger.Fatal("server initialization failed", log.Error(err))
	}
}
