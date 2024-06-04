package server

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/config"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/router"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/storage"
)

// Init starts the application server.
func Init(ctx context.Context, cfg *config.Config, logger *log.Logger) error {
	logger.Debug("initializing server")

	redisPrimaryClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})
	redisReaderClient := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	defer func() {
		if err := redisPrimaryClient.Close(); err != nil {
			logger.Error("failed to close redis primary connection pool", log.Error(err))
		} else {
			logger.Debug("redis primary connection closed")
		}
		if err := redisReaderClient.Close(); err != nil {
			logger.Error("failed to close redis reader connection pool", log.Error(err))
		} else {
			logger.Debug("redis reader connection closed")
		}
	}()

	logger.Debug("redis connected")

	refreshTokenStorage := storage.NewRedisRepository("[Subject2RefreshToken]", redisPrimaryClient, redisReaderClient)
	accessTokenStorage := storage.NewRedisRepository("[RefreshToken2AccessToken]", redisPrimaryClient, redisReaderClient)
	flowRepository := storage.NewRedisRepository("[Subject2Flow]", redisPrimaryClient, redisReaderClient)

	privateKeyBytes := []byte(cfg.AuthenticationService.PrivateKey)
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return err
	}

	authClient := domain.NewClient(
		logger,
		domain.NewUserServiceClient(cfg.UserService.Timeout*time.Second, cfg.UserService.Endpoint),
		refreshTokenStorage,
		accessTokenStorage,
		flowRepository,
		cfg.AuthenticationService.Claims.Issuer,
		map[string]time.Duration{
			domain.AccessTokenExpiration:       time.Duration(cfg.AuthenticationService.Claims.AccessTokenExpiration) * time.Second,
			domain.AnonymousExpiration:         time.Duration(cfg.AuthenticationService.Claims.AnonymousExpiration) * time.Second,
			domain.ShortRefreshTokenExpiration: time.Duration(cfg.AuthenticationService.Claims.ShortRefreshTokenExpiration) * time.Second,
			domain.LongRefreshTokenExpiration:  time.Duration(cfg.AuthenticationService.Claims.LongRefreshTokenExpiration) * time.Second,
		},
		privateKey,
	)

	server := http.Server{
		Addr:         cfg.AuthenticationService.Server.Port,
		Handler:      router.New(authClient, logger).GetHandler(),
		ReadTimeout:  cfg.AuthenticationService.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.AuthenticationService.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.AuthenticationService.Server.IdleTimeout * time.Second,
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
