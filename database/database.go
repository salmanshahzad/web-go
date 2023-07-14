package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Ctx context.Context = context.Background()
	Db  *gorm.DB
	Rdb *redis.Client
)
