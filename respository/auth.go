package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

func Register(ctx *gin.Context, pool *pgxpool.Pool) models.User {
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

	_, err = pool.Exec(context.Background(), "INSERT INTO users (fullname, email, password, role, profileId, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)", input.Fullname, input.Email, hash, input.Role, input.Profileid, now, now)

	if err != nil {
		fmt.Println("Error insert user:", err)
	}

	input.Password = string(hash)
	input.Created_at = &now
	input.Updated_at = &now

	return input
}

func FindUserEmail(pool *pgxpool.Pool, email string) (models.User, error) {
	var users models.User

	row := pool.QueryRow(context.Background(), "SELECT id, fullname, email, password, role, profileId, created_at, updated_at FROM users WHERE email = $1", email)

	err := row.Scan(&users.Id, &users.Fullname, &users.Email, &users.Password, &users.Role, &users.Profileid, &users.Created_at, &users.Updated_at)

	users.Password = strings.TrimSpace(users.Password)

	return users, err
}

func VerifPassword(inputPassword string, hashPassword string) bool {
	ok, err := argon2.VerifyEncoded([]byte(hashPassword), []byte(inputPassword))

	if err != nil {
		fmt.Println("Error : Password not metmatch")
	}

	return ok
}

func UpdatePassword(pool *pgxpool.Pool, email string, newPassword string) error{
	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(newPassword))
	if err != nil {
		fmt.Println("Failed to hash password", err)

	}

	_, err = pool.Exec(context.Background(), "UPDATE users SET password = $1, updated_at=$2 WHERE email=$3", hash, time.Now(), email)

	if err != nil {
		fmt.Println("Error updating password:", err)
	}
	
	return nil
}
