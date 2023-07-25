package database

import (
	"github.com/redis/go-redis/v9"
)

var (
	Db  *Queries
	Rdb *redis.Client
)
