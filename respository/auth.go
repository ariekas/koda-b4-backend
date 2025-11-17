package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matthewhartstonge/argon2"
)

func Register(ctx *gin.Context, pool *pgxpool.Pool) (models.User, error) {
	var input models.RegisterRequest
	var checkEmail bool

	if err := ctx.BindJSON(&input); err != nil {
		return models.User{}, fmt.Errorf("invalid request body")
	}

	if input.Fullname == "" {
		return models.User{}, fmt.Errorf("fullname is required")
	}
	if input.Email == "" {
		return models.User{}, fmt.Errorf("email is required")
	}
	if input.Password == "" {
		return models.User{}, fmt.Errorf("password is required")
	}
	if len(input.Password) <= 6 {
		return models.User{}, fmt.Errorf("password must be more than 6 characters")
	}

	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(input.Password))
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password")
	}

	err = pool.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`,
		input.Email,
	).Scan(&checkEmail)

	if err != nil {
		return models.User{}, fmt.Errorf("error checking email")
	}
	if checkEmail {
		return models.User{}, fmt.Errorf("email already registered")
	}

	role := "user"

	now := time.Now()
	var profileID *int = nil

	_, err = pool.Exec(
		context.Background(),
		`INSERT INTO users (fullname, email, password, role, profile_id, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		input.Fullname, input.Email, hash, role, profileID, now, now,
	)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to insert user")
	}

	user := models.User{
		Fullname:  input.Fullname,
		Email:     input.Email,
		Password:  string(hash),
		Role:      role,
		ProfileID: profileID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	return user, nil
}


func FindUserEmail(pool *pgxpool.Pool, email string) (models.User, error) {
	var user models.User

	row := pool.QueryRow(context.Background(),`
		SELECT u.id, u.fullname, u.email, u.password, u.role, u.profile_id,
		       p.pic, p.phone, p.address, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profile p ON u.profile_id = p.id
		WHERE u.email = $1
	`,email)

	err := row.Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.ProfileID,
		&user.Pic,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, fmt.Errorf("no user found with this email address, %w", err)
	}
	return user, nil
}

func VerifPassword(inputPassword string, hashPassword string) bool {
	ok, err := argon2.VerifyEncoded([]byte(hashPassword), []byte(inputPassword))

	if err != nil {
		fmt.Println("Error : Password not metmatch, ",err)
	}

	return ok
}

func UpdatePassword(pool *pgxpool.Pool, email string, newPassword string) error {

	if len(newPassword) <= 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}

	argon := argon2.DefaultConfig()
	hash, err := argon.HashEncoded([]byte(newPassword))
	if err != nil {
		fmt.Println("Failed to hash password", err)

	}

	_, err = pool.Exec(context.Background(), "UPDATE users SET password = $1, updated_at=$2 WHERE email=$3", hash, time.Now(), email)

	if err != nil {
		return fmt.Errorf("failed updating password, %w", err)
	}

	return nil
}
