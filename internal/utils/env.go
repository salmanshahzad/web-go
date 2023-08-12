package utils

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Environment struct {
	CorsOrigins   string `env:"CORS_ORIGINS,default=*"`
	DbHost        string `env:"DB_HOST,required"`
	DbName        string `env:"DB_NAME,required"`
	DbPassword    string `env:"DB_PASSWORD,required"`
	DbPort        int    `env:"DB_PORT,default=5432"`
	DbUser        string `env:"DB_USER,required"`
	Port          int    `env:"PORT,default=1024"`
	RedisHost     string `env:"REDIS_HOST,required"`
	RedisPassword string `env:"REDIS_PASSWORD,required"`
	RedisPort     int    `env:"REDIS_PORT,default=6379"`
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
