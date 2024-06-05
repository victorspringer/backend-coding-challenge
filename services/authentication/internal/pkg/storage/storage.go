package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/domain"
)

var errRedisNotFound = "redis: nil"

type redisRepository struct {
	redisPrefix        string
	redisPrimaryClient *redis.Client
	redisReaderClient  *redis.Client
}

// NewRedisRepository returns a new instance of the Redis repository.
func NewRedisRepository(redisPrefix string, redisPrimaryClient, redisReaderClient *redis.Client) domain.Repository {
	return &redisRepository{
		redisPrefix:        redisPrefix,
		redisPrimaryClient: redisPrimaryClient,
		redisReaderClient:  redisReaderClient,
	}
}

// Keys retrieves keys matching the specified pattern.
func (r *redisRepository) Keys(ctx context.Context, pattern string) ([]string, error) {
	var (
		keys   []string = make([]string, 0)
		kk     []string
		cursor uint64
		err    error
	)

	for {
		kk, cursor, err = r.redisReaderClient.Scan(ctx, cursor, `\`+r.redisPrefix+pattern+`-*`, 0).Result()
		if err != nil {
			if err.Error() == errRedisNotFound {
				return []string{}, fmt.Errorf("no keys found for pattern %s", pattern)
			}
			return []string{}, err
		}
		keys = append(keys, kk...)
		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return []string{}, fmt.Errorf("no keys found for pattern %s", pattern)
	}

	for i := range keys {
		keys[i] = strings.Replace(keys[i], r.redisPrefix, "", 1)
	}

	return keys, nil
}

// Get retrieves the value associated with the specified key.
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.redisReaderClient.Get(ctx, r.redisPrefix+key).Result()
	if err != nil {
		if err.Error() == errRedisNotFound {
			return "", fmt.Errorf("value not found for key %s", key)
		}
		return "", err
	}
	if value == "" {
		return "", fmt.Errorf("value not found for key %s", key)
	}
	return value, nil
}

// Set sets the value associated with the specified key with an optional expiration duration.
func (r *redisRepository) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	return r.redisPrimaryClient.Set(ctx, r.redisPrefix+key, value, expiration).Err()
}

// Del deletes the value associated with the specified key.
func (r *redisRepository) Del(ctx context.Context, key string) error {
	err := r.redisPrimaryClient.Del(ctx, r.redisPrefix+key).Err()
	if err != nil {
		return err
	}
	return nil
}
