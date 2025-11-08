package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/matthewhartstonge/argon2"
)

func Register(ctx *gin.Context, conn *pgx.Conn) models.User {
	var input models.User
	argon := argon2.DefaultConfig()

	err := ctx.BindJSON(&input)
	if err != nil {
		fmt.Println("Error: Invalid type much json")
	}

	hash, err := argon.HashEncoded([]byte(input.Password))
	if err != nil {
		fmt.Println("Error : Failed to hash password")
	}

	now := time.Now()

	_, err = conn.Exec(context.Background(), "INSERT INTO users (fullname, email, password, role, profileId, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", input.Fullname, input.Email, hash, input.Role, input.Profileid, now, now)

	if err != nil {
		fmt.Println("Error insert user:", err)
	}

	input.Password = string(hash)
	input.Created_at = now
	input.Updated_at = now

	
	return input
}