package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateVariant(pool *pgxpool.Pool, req *models.Variant) (*models.Variant, error) {

	now := time.Now()
	var variant models.Variant

	err := pool.QueryRow(context.Background(), `
		INSERT INTO variant_products (name, additional_costs, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, additional_costs, created_at, updated_at
	`, req.Name, req.AdditionalCosts, now, now).Scan(
		&variant.ID,
		&variant.Name,
		&variant.AdditionalCosts,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		return &models.Variant{}, fmt.Errorf("failed to create variant product: %v", err)
	}

	return &variant, nil
}

func ListVariant(pool *pgxpool.Pool) ([]models.Variant, error) {

	rows, err := pool.Query(context.Background(), `
		SELECT id, name, additional_costs, created_at, updated_at
		FROM variant_products
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all variant products: %v", err)
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		var variant models.Variant
		err := rows.Scan(
			&variant.ID,
			&variant.Name,
			&variant.AdditionalCosts,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant product: %v", err)
		}
		variants = append(variants, variant)
	}

	return variants, nil
}

func DetailVariant(pool *pgxpool.Pool, id int) (*models.Variant, error) {

	var variant models.Variant
	err := pool.QueryRow(context.Background(), `
		SELECT id, name, additional_costs, created_at, updated_at
		FROM variant_products
		WHERE id = $1
	`, id).Scan(
		&variant.ID,
		&variant.Name,
		&variant.AdditionalCosts,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get variant product: %v", err)
	}

	return &variant, nil
}

func UpdateVariant(pool *pgxpool.Pool, id int, req *models.Variant) (*models.Variant, error) {

	var variant models.Variant
	err := pool.QueryRow(context.Background(), `
		UPDATE variant_products
		SET name = $1, additional_costs = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, name, additional_costs, created_at, updated_at
	`, req.Name, req.AdditionalCosts, time.Now(), id).Scan(
		&variant.ID,
		&variant.Name,
		&variant.AdditionalCosts,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update variant product: %v", err)
	}

	return &variant, nil
}

func DeleteVariant(pool *pgxpool.Pool, id int) error {

	result, err := pool.Exec(context.Background(), `DELETE FROM variant_products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete variant product: %v", err)
	}

	rowsAffected := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("variant product not found")
	}

	return nil
}
