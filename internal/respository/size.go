package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateSize(req *models.Size, pool *pgxpool.Pool) (*models.Size, error) {

	size := &models.Size{}
	err := pool.QueryRow(context.Background(), `INSERT INTO size_products (name, additional_costs, created_at, updated_at) 
		VALUES ($1, $2, NOW(), NOW()) 
		RETURNING id, name, additional_costs, created_at, updated_at`, req.Name, req.AdditionalCosts).Scan(
		&size.Id,
		&size.Name,
		&size.AdditionalCosts,
		&size.CreatedAt,
		&size.UpdatedAt,
	)

	if err != nil {
		fmt.Println("failed to create size: %w", err)
	}

	return size, nil
}

func ListSize(pool *pgxpool.Pool) ([]models.Size, error) {

	rows, err := pool.Query(context.Background(), `SELECT id, name, additional_costs, created_at, updated_at 
		FROM size_products 
		ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("failed to get sizes: %w", err)
	}
	defer rows.Close()

	sizes := []models.Size{}
	for rows.Next() {
		var size models.Size
		err := rows.Scan(
			&size.Id,
			&size.Name,
			&size.AdditionalCosts,
			&size.CreatedAt,
			&size.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan size: %w", err)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

func DetailSize(pool *pgxpool.Pool, id int) (*models.Size, error) {
	size := &models.Size{}
	err := pool.QueryRow(context.Background(), `SELECT id, name, additional_costs, created_at, updated_at 
		FROM size_products 
		WHERE id = $1`, id).Scan(
		&size.Id,
		&size.Name,
		&size.AdditionalCosts,
		&size.CreatedAt,
		&size.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get size: %w", err)
	}

	return size, nil
}

func UpdateSize(pool *pgxpool.Pool, id int, req *models.Size) (*models.Size, error) {

	size := &models.Size{}
	err := pool.QueryRow(context.Background(), `
		UPDATE size_products 
		SET name = $1, additional_costs = $2, updated_at = NOW() 
		WHERE id = $3 
		RETURNING id, name, additional_costs, created_at, updated_at
	`, req.Name, req.AdditionalCosts, id).Scan(
		&size.Id,
		&size.Name,
		&size.AdditionalCosts,
		&size.CreatedAt,
		&size.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update size: %w", err)
	}

	return size, nil
}

func DeleteSize(pool *pgxpool.Pool, id int) error {

	result, err := pool.Exec(context.Background(), `DELETE FROM size_products WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete size: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("size with id %d not found", id)
	}

	return nil
}
