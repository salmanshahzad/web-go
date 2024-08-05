package utils

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}
