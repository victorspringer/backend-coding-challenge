package domain

import (
	"context"
	"time"
)

// Repository is the interface for the authenticator storage repository.
type Repository interface {
	// Keys retrieves keys matching the specified pattern.
	Keys(ctx context.Context, pattern string) ([]string, error)
	// Get retrieves the value associated with the specified key.
	Get(ctx context.Context, key string) (string, error)
	// Set sets the value associated with the specified key with an optional expiration duration.
	Set(ctx context.Context, key, value string, expiration time.Duration) error
	// Del deletes the value associated with the specified key.
	Del(ctx context.Context, key string) error
}
