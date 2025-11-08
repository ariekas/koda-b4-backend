package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetDataUsers(conn *pgx.Conn) ([]models.User, error) {
	var dataUser []models.User
	rows, err := conn.Query(context.Background(), `SELECT id, fullname, email,password, role, profileid, created_at, updated_at FROM users`)
	if err != nil {
		fmt.Println("Error: Failed get data users")
	}

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.Id,
			&user.Fullname,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Profileid,
			&user.Created_at,
			&user.Updated_at,
		)
		if err != nil {
			fmt.Println("Error scanning user:", err)
		}
		dataUser = append(dataUser, user)
	}

	return dataUser, nil
}
