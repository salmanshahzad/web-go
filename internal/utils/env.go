package utils

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Environment struct {
	CorsOrigins string `env:"CORS_ORIGINS,default=*"`
	DatabaseUrl string `env:"DATABASE_URL,required"`
	Port        int    `env:"PORT,default=1024"`
	RedisUrl    string `env:"REDIS_URL,required"`
}

func InitEnv() (*Environment, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	ctx := context.Background()
	env := new(Environment)
	if err := envconfig.Process(ctx, env); err != nil {
		return nil, err
	}

	return env, nil
}
