package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDiscount(pool *pgxpool.Pool) ([]models.Discount, error) {
	var discounts []models.Discount

	rows, err := pool.Query(context.Background(),
		`SELECT id, name, discount_percentage, created_at, updated_at 
		 FROM discounts 
		 ORDER BY id DESC`)
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

func CreateDiscount(pool *pgxpool.Pool, input models.Discount) (*models.Discount, error) {
	now := time.Now()

	var discount models.Discount

	err := pool.QueryRow(context.Background(),
		`INSERT INTO discounts (name, discount_percentage, created_at, updated_at)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, name, discount_percentage, created_at, updated_at`,
		input.Name,
		input.Diskon,
		now,
		now,
	).Scan(
		&discount.Id,
		&discount.Name,
		&discount.Diskon,
		&discount.CreatedAt,
		&discount.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create discount: %w", err)
	}

	return &discount, nil
}

func DetailDiscount(pool *pgxpool.Pool, id int) (*models.Discount, error) {
	var d models.Discount

	err := pool.QueryRow(context.Background(),
		`SELECT id, name, discount_percentage, created_at, updated_at
		 FROM discounts 
		 WHERE id = $1`,
		id,
	).Scan(
		&d.Id,
		&d.Name,
		&d.Diskon,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed get detail discount: %w", err)
	}

	return &d, nil
}

func UpdateDiscount(pool *pgxpool.Pool, id int, input models.Discount) (*models.Discount, error) {
	now := time.Now()

	var updated models.Discount

	err := pool.QueryRow(context.Background(),
		`UPDATE discounts 
		 SET name = $1, discount_percentage = $2, updated_at = $3
		 WHERE id = $4
		 RETURNING id, name, discount_percentage, created_at, updated_at`,
		input.Name,
		input.Diskon,
		now,
		id,
	).Scan(
		&updated.Id,
		&updated.Name,
		&updated.Diskon,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update discount: %w", err)
	}

	return &updated, nil
}

func DeleteDiscount(pool *pgxpool.Pool, id int) error {
	_, err := pool.Exec(context.Background(),
		`DELETE FROM discounts WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete discount: %w", err)
	}

	return nil
}
