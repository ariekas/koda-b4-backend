package controller

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB() *pgx.Conn{
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		panic("Error : Failed to connect database")
	}

	return conn
}