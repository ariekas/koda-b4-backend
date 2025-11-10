package controller

import (
	"back-end-coffeShop/lib/config"
	"context"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() *pgx.Conn{
	dbURL:= config.ReadEnvDb()

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		panic("Error : Failed to connect database")
	}

	return conn
}