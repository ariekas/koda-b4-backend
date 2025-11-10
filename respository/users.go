package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDataUsers(pool *pgxpool.Pool) ([]models.User, error) {
	var dataUser []models.User
	rows, err := pool.Query(context.Background(), `SELECT 
  u.id AS user_id,
  u.fullname,
  u.email,
  p.pic,
  p.phone,
  p.address
FROM users u
LEFT JOIN profile p ON u.profileid = p.id
ORDER BY u.id;
`)
	if err != nil {
		fmt.Println("Error: Failed get data users", err)
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.Fullname,
			&user.Email,
			&user.Pic,
			&user.Phone,	
			&user.Address,
		)
		if err != nil {
			fmt.Println("Error scanning user:", err)
		}
		dataUser = append(dataUser, user)
	}

	return dataUser, nil
}

func DeleteUser(pool *pgxpool.Pool, ctx *gin.Context) error {
	id := ctx.Param("id")
	_, err := pool.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)

	return err
}

func UpdateRole(pool *pgxpool.Pool, ctx *gin.Context, userId int, newRole string) error {
	_, err := pool.Exec(context.Background(), "UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2", newRole, userId)

	if err != nil {
		fmt.Println("Error: Failed to update role user ",err)
	}

	return nil
}