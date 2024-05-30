package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/database"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/router"
)

// Init starts the application server.
func Init(cfg *config.Config, logger *log.Logger) error {
	db, err := database.New()
	if err != nil {
		logger.Fatal("failed to init database", log.Error(err))
	}
	// defer db.Close()

	server := http.Server{
		Addr:         cfg.UserService.Server.Port,
		Handler:      router.New(db, logger).GetHandler(),
		ReadTimeout:  cfg.UserService.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.UserService.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.UserService.Server.IdleTimeout * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("ListenAndServe failed", log.Error(err))
		}
	}()

	logger.Info("server is listening", log.String("addr", server.Addr))

	// channel to listen for interrupt or terminate signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// block until a signal is received
	<-stop

	// context with a timeout to enable graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), server.WriteTimeout*time.Second)
	defer cancel()

	// gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server shutdown failed")
	} else {
		logger.Info("server gracefully stopped")
	}

	return nil
}
