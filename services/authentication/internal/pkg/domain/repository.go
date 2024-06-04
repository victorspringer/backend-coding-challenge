package domain

import (
	"context"
	"time"
)

// Repository is the interface for the authenticator storage.
type Repository interface {
	Keys(ctx context.Context, pattern string) ([]string, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}
