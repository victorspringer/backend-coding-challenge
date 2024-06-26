package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	authClient "github.com/victorspringer/backend-coding-challenge/services/authentication/pkg/client"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/database"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/router"
)

// Init starts the application server.
func Init(ctx context.Context, cfg *config.Config, logger *log.Logger) error {
	logger.Debug("initializing server")

	db, err := database.New(
		ctx,
		logger,
		cfg.MongoDB.URI,
		cfg.MongoDB.DBName,
		cfg.MongoDB.Collection,
		cfg.MongoDB.Timeout*time.Second,
	)
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}

	defer func() {
		if err := db.Close(ctx); err != nil {
			logger.Error("failed to close database connection pool", log.Error(err))
		} else {
			logger.Debug("database connection closed")
		}
	}()

	logger.Debug("database connected")

	ac := authClient.NewClient(
		cfg.AuthenticationService.URL,
		cfg.AuthenticationService.Timeout*time.Second,
		logger,
	)

	server := http.Server{
		Addr:         cfg.RatingService.Server.Port,
		Handler:      router.New(db, logger, ac).GetHandler(),
		ReadTimeout:  cfg.RatingService.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.RatingService.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.RatingService.Server.IdleTimeout * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server failed to listen and serve", log.Error(err))
		}
	}()

	logger.Info("server is listening", log.String("addr", server.Addr))

	// channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// block until a signal is received
	<-stop

	// context with a timeout to enable graceful shutdown
	ctx, cancel := context.WithTimeout(ctx, server.WriteTimeout*time.Second)
	defer cancel()

	// gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server shutdown failed")
	} else {
		logger.Info("server gracefully stopped")
	}

	return nil
}
