package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateDelivery(pool *pgxpool.Pool, delivery *models.Deliverys) error {
	now := time.Now()
	err := pool.QueryRow(context.Background(), `
		INSERT INTO deliverys ( name, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`, delivery.Name, delivery.Price, now, now).
		Scan(&delivery.ID, &delivery.CreatedAt, &delivery.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create delivery: %w", err)
	}

	return nil
}

func ListDelivery(pool *pgxpool.Pool) ([]models.Deliverys, error) {
	rows, err := pool.Query(context.Background(), `
		SELECT id, name, price, created_at, updated_at
		FROM deliverys
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []models.Deliverys
	for rows.Next() {
		var delivery models.Deliverys
		err := rows.Scan(&delivery.ID, &delivery.Name, &delivery.Price, &delivery.CreatedAt, &delivery.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %w", err)
		}
		deliveries = append(deliveries, delivery)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating deliveries: %w", err)
	}

	return deliveries, nil
}

func DetailDelivery(pool *pgxpool.Pool, id int) (*models.Deliverys, error) {
	var delivery models.Deliverys
	err := pool.QueryRow(context.Background(), `
		SELECT id, name, price, created_at, updated_at
		FROM deliverys
		WHERE id = $1
	`, id).
		Scan(&delivery.ID, &delivery.Name, &delivery.Price, &delivery.CreatedAt, &delivery.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}

	return &delivery, nil
}

func UpdateDelivery(pool *pgxpool.Pool, delivery *models.Deliverys) error {
	now := time.Now()
	err := pool.QueryRow(context.Background(), `
		UPDATE deliverys
		SET name = $2, price = $3, updated_at = $4
		WHERE id = $1
		RETURNING updated_at
	`, delivery.ID, delivery.Name, delivery.Price, now).
		Scan(&delivery.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update delivery: %w", err)
	}

	return nil
}

func DeleteDelivery(pool *pgxpool.Pool, id int) error {
	cmdTag, err := pool.Exec(context.Background(), `DELETE FROM deliverys WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete delivery: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("delivery with id %d not found", id)
	}

	return nil
}
