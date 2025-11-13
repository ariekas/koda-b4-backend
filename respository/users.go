package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)
type PaginationResponseUser struct {
	Data       []models.User   `json:"data"`
	Page       int             `json:"page"`
	Limit      int             `json:"limit"`
	Total      int             `json:"total"`
	TotalPages int             `json:"total_pages"`
	Links      map[string]string `json:"links"`
}

func GetDataUsers(pool *pgxpool.Pool, page int) (PaginationResponseUser, error) {
	var dataUser []models.User
	limit := 50
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		fmt.Println("Error counting users:", err)
	}

	rows, err := pool.Query(
		context.Background(),
		`SELECT u.id, u.fullname, u.email, u.role, u.profile_id,
		        p.pic, p.phone, p.address, u.created_at, u.updated_at
		 FROM users u
		 LEFT JOIN profile p ON u.profile_id = p.id
		 ORDER BY u.id
		 OFFSET $1 LIMIT $2`,
		offset, limit,
	)
	if err != nil {
		fmt.Println("Error: Failed get data users", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.Id,
			&user.Fullname,
			&user.Email,
			&user.Role,
			&user.ProfileID,
			&user.Pic,
			&user.Phone,
			&user.Address,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			fmt.Println("Error scanning user:", err)
			continue
		}
		dataUser = append(dataUser, user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := make(map[string]string)
	if page > 1 {
		links["prev"] = fmt.Sprintf("/users?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/users?page=%d", page+1)
	}

	return PaginationResponseUser{
		Data:       dataUser,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}

func DeleteUser(pool *pgxpool.Pool, userId int) error {
	_, err := pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", userId)
	return err
}

func UpdateRole(pool *pgxpool.Pool, userId int, newRole string) error {
	_, err := pool.Exec(context.Background(), "UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2", newRole, userId)
	return err
}

func UpdateProfile(pool *pgxpool.Pool, userId int, input models.UpdateProfileRequest) error {
	tx, err := pool.Begin(context.Background())
	if err != nil {
		fmt.Println("failed to start transaction", err)
	}
	defer tx.Rollback(context.Background())

	var profileId int
	err = tx.QueryRow(context.Background(), `
		SELECT profile_id FROM users WHERE id = $1
	`, userId).Scan(&profileId)
	if err != nil {
		fmt.Println("user not found or profile not assigned", err)
	}

	_, err = tx.Exec(context.Background(), `
		UPDATE profile
		SET pic = $1, phone = $2, address = $3, updated_at = NOW()
		WHERE id = $4
	`, input.Pic, input.Phone, input.Address, profileId)
	if err != nil {
		fmt.Println("failed to update profile", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		fmt.Println("failed to commit transaction", err)
	}

	return nil
}

func GetUserById(pool *pgxpool.Pool, userId int) (models.User, error) {
	var user models.User
	err := pool.QueryRow(context.Background(), `
		SELECT 
			u.id, u.fullname, u.email, u.role, u.profile_id,
			p.pic, p.phone, p.address, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profile p ON u.profile_id = p.id
		WHERE u.id = $1
	`, userId).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Role,
		&user.ProfileID,
		&user.Pic,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		fmt.Println("failed to get user profile", err)
	}
	return user, nil
}