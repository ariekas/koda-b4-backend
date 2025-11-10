package controller

import (
	"back-end-coffeShop/lib/config"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() *pgxpool.Pool{
	dbURL:= config.ReadEnvDb()

	pool, err := pgxpool.New(context.Background(), dbURL)

	if err != nil {
		panic("Error : Failed to connect database")
	}

	return pool
}