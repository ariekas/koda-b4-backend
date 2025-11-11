package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func Redis() *redis.Client {
	godotenv.Load()
	redisUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	rdb := redis.NewClient(&redis.Options{
		Addr: redisUrl,
		Password: password,
		DB: 0,
	})

	return rdb
}