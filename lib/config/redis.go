package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
    _ = godotenv.Load()

    opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
    if err != nil {
        panic(err)
    }

    return redis.NewClient(opt)
}
