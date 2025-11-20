package respository

import (
	"back-end-coffeShop/lib/config"
	"back-end-coffeShop/models"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	result, err := pool.Exec(context.Background(),
		"DELETE FROM users WHERE id = $1",
		userId,
	)

	if err != nil {
		return fmt.Errorf("failed to delete user, %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}


func UpdateRole(pool *pgxpool.Pool, userId int, newRole string) error {
	if strings.TrimSpace(newRole) == "" {
		return fmt.Errorf("role is required")
	}

	result, err := pool.Exec(context.Background(),
		"UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2",
		newRole, userId,
	)
	if err != nil {
		return fmt.Errorf("failed to update role: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func UpdateProfile(pool *pgxpool.Pool, userId int, input models.UpdateProfileRequest) error {
	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	var profileId *int
	err = tx.QueryRow(ctx, `
		SELECT profile_id FROM users WHERE id = $1
	`, userId).Scan(&profileId)

	if err != nil {
		return fmt.Errorf("failed to get user profile: %v", err)
	}

	if profileId == nil {
		var newProfileId int

		err = tx.QueryRow(ctx, `
			INSERT INTO profile (pic, phone, address, created_at, updated_at)
			VALUES ('', '', '', NOW(), NOW())
			RETURNING id
		`).Scan(&newProfileId)

		if err != nil {
			return fmt.Errorf("failed to create new profile: %v", err)
		}

		_, err = tx.Exec(ctx, `
			UPDATE users SET profile_id = $1 WHERE id = $2
		`, newProfileId, userId)

		if err != nil {
			return fmt.Errorf("failed to assign new profile to user: %v", err)
		}

		profileId = &newProfileId
	}

	var fileName *string = nil

	if input.PicFile != nil {
		ext := strings.ToLower(filepath.Ext(input.PicFile.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return fmt.Errorf("unsupported image format (only jpg, jpeg, png)")
		}

		useCloud := os.Getenv("CLOUDINARY_API_KEY") != "" &&
			os.Getenv("CLOUDINARY_API_SECRET") != "" &&
			os.Getenv("CLOUDINARY_NAME") != ""

		newName := fmt.Sprintf("profile_%d_%d%s", userId, time.Now().Unix(), ext)

		if useCloud {
			uploadedURL, err := config.UploaderFile(input.PicFile, "pic", newName)
			if err != nil {
				return fmt.Errorf("failed upload cloudinary: %v", err)
			}

			fileName = &uploadedURL

		} else {
			savePath := "uploads/profile/" + newName

			if err := os.MkdirAll("uploads/profile", os.ModePerm); err != nil {
				return fmt.Errorf("failed to create folder: %v", err)
			}

			if err := SaveUploadedFile(input.PicFile, savePath); err != nil {
				return fmt.Errorf("failed to save local file: %v", err)
			}

			fileName = &newName
		}
	}
	_, err = tx.Exec(ctx, `
		UPDATE profile
		SET 
			pic = COALESCE($1, pic),
			phone = COALESCE($2, phone),
			address = COALESCE($3, address),
			updated_at = NOW()
		WHERE id = $4
	`, fileName, input.Phone, input.Address, *profileId)

	if err != nil {
		return fmt.Errorf("failed to update profile: %v", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}


func GetUserByToken(pool *pgxpool.Pool, userId int) (models.User, error) {
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
		return models.User{}, fmt.Errorf("failed to get user profile, %w", err)
	}

	return user, nil
}