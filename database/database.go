package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx context.Context = context.Background()
	Db  *Queries
	Rdb *redis.Client
)
