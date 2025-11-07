package controller

import (
	"back-end-coffeShop/models"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDB(ctx *gin.Context) *pgx.Conn{
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbURL)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error : Failed connect to database",
		})
	}

	return conn
}