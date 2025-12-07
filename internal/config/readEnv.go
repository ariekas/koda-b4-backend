package config

import (
	"os"

	"github.com/joho/godotenv"
)

func ReadENV() string {
	godotenv.Load()
	JWTtoken  := os.Getenv("JWT_TOKEN")

	return JWTtoken
}

func ReadEnvDb() string {
	godotenv.Load()

	DbUrl := os.Getenv("DATABASE_URL")

	return DbUrl
}

func ReadEnvUrl() string{
	godotenv.Load()

	url := os.Getenv("ORIGIN_URL")

	return url
}