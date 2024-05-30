package main

import (
	"os"

	"github.com/victorspringer/backend-coding-challenge/services/user/internal/app/server"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
)

// TODO: implement repository
// tests

func main() {
	env := os.Getenv("ENVIRONMENT")
	cfg, err := config.New(env)
	if err != nil {
		panic(err)
	}

	logger := log.New(cfg.UserService.LogLevel)

	if err := server.Init(cfg, logger); err != nil {
		logger.Fatal("failed to initialize server", log.Error(err))
	}
}
