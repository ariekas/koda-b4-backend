package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTax(pool *pgxpool.Pool, tax *models.CreateTaxRequest) error {
	now := time.Now()

	var id int
	var createdAt, updatedAt time.Time

	err := pool.QueryRow(context.Background(), `
		INSERT INTO taxes (name, tax, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`,
		tax.Name,
		tax.Tax,
		now,
		now,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return fmt.Errorf("failed to create tax: %w", err)
	}

	return nil
}

func ListTax(pool *pgxpool.Pool) ([]models.Taxe, error) {

	rows, err := pool.Query(context.Background(), `
		SELECT id, name, tax, created_at, updated_at
		FROM taxes
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get taxes: %w", err)
	}
	defer rows.Close()

	var taxes []models.Taxe
	for rows.Next() {
		var tax models.Taxe
		err := rows.Scan(&tax.ID, &tax.Name, &tax.Tax, &tax.CreatedAt, &tax.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tax: %w", err)
		}
		taxes = append(taxes, tax)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating taxes: %w", err)
	}

	return taxes, nil
}

func DetailTax(pool *pgxpool.Pool, id int) (*models.Taxe, error) {

	var tax models.Taxe
	err := pool.QueryRow(context.Background(), `
		SELECT id, name, tax, created_at, updated_at
		FROM taxes
		WHERE id = $1
	`, id).
		Scan(&tax.ID, &tax.Name, &tax.Tax, &tax.CreatedAt, &tax.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get tax: %w", err)
	}

	return &tax, nil
}

func UpdateTax(pool *pgxpool.Pool, tax *models.Taxe) error {

	now := time.Now()
	err := pool.QueryRow(context.Background(), `
		UPDATE taxes
		SET name = $2, tax = $3, updated_at = $4
		WHERE id = $1
		RETURNING updated_at
	`, tax.ID, tax.Name, tax.Tax, now).
		Scan(&tax.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update tax: %w", err)
	}

	return nil
}

func DeleteTax(pool *pgxpool.Pool, id int) error {

	cmdTag, err := pool.Exec(context.Background(), `DELETE FROM taxes WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete tax: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("tax with id %d not found", id)
	}

	return nil
}
