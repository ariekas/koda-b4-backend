package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDiscount(pool *pgxpool.Pool) ([]models.Discount, error) {
	var discounts []models.Discount

	rows, err := pool.Query(context.Background(),
		`SELECT id, name, diskon, created_at, updated_at FROM discounts ORDER BY id DESC`)
	if err != nil {
		return nil, fmt.Errorf("error query discounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var d models.Discount
		err := rows.Scan(
			&d.Id,
			&d.Name,
			&d.Diskon,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scan discounts: %w", err)
		}

		discounts = append(discounts, d)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return discounts, nil
}

