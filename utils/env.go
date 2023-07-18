package utils

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Environment struct {
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

var Env Environment

func InitEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, &Env); err != nil {
		return err
	}

	return nil
}
