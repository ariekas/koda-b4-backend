package respository

import (
	"back-end-coffeShop/internal/models"
	
	"context"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetPaymentMethods(pool *pgxpool.Pool, page int) (models.PaginationResponse, error) {
	var paymentMethods []models.PaymentMethods
	limit := 10
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM payment_methods").Scan(&total)
	if err != nil {
		return models.PaginationResponse{}, fmt.Errorf("error counting payment methods: %w", err)
	}

	rows, err := pool.Query(context.Background(), 
		"SELECT id, name, image_payment, created_at, updated_at FROM payment_methods ORDER BY created_at DESC OFFSET $1 LIMIT $2", 
		offset, limit)
	if err != nil {
		return models.PaginationResponse{}, fmt.Errorf("failed to get payment methods: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pm models.PaymentMethods
		if err := rows.Scan(&pm.Id, &pm.Name, &pm.ImagePayment, &pm.CreatedAt, &pm.UpdatedAt); err != nil {
			return models.PaginationResponse{}, fmt.Errorf("error scanning payment method: %w", err)
		}
		paymentMethods = append(paymentMethods, pm)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := map[string]string{}
	if page > 1 {
		links["prev"] = fmt.Sprintf("/payment-methods?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/payment-methods?page=%d", page+1)
	} else {
		links["next"] = "null"
	}

	return models.PaginationResponse{
		Data:       paymentMethods,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}

func CreatePaymentMethod(pool *pgxpool.Pool, input models.PaymentMethods) (models.PaymentMethods, error) {
	err := pool.QueryRow(context.Background(),
		"INSERT INTO payment_methods (name, image_payment) VALUES ($1, $2) RETURNING id, created_at, updated_at",
		input.Name, input.ImagePayment,
	).Scan(&input.Id, &input.CreatedAt, &input.UpdatedAt)

	if err != nil {
		return models.PaymentMethods{}, fmt.Errorf("failed to insert payment method: %w", err)
	}

	return input, nil
}

func GetPaymentMethodById(pool *pgxpool.Pool, id int) (models.PaymentMethods, error) {
	var pm models.PaymentMethods
	err := pool.QueryRow(context.Background(),
		"SELECT id, name, image_payment, created_at, updated_at FROM payment_methods WHERE id=$1", id,
	).Scan(&pm.Id, &pm.Name, &pm.ImagePayment, &pm.CreatedAt, &pm.UpdatedAt)

	if err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error: Payment method not found: %w", err)
	}

	return pm, nil
}

func EditPaymentMethod(pool *pgxpool.Pool, id int, input models.PaymentMethods) (models.PaymentMethods, error) {
	pm, err := GetPaymentMethodById(pool, id)
	if err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error: not found payment method: %w", err)
	}

	if input.Name != "" {
		pm.Name = input.Name
	}
	if input.ImagePayment != "" {
		pm.ImagePayment = input.ImagePayment
	}

	_, err = pool.Exec(context.Background(),
		"UPDATE payment_methods SET name=$1, image_payment=$2, updated_at=NOW() WHERE id=$3",
		pm.Name, pm.ImagePayment, id,
	)

	if err != nil {
		return models.PaymentMethods{}, fmt.Errorf("error: failed to update payment method: %w", err)
	}

	pm, err = GetPaymentMethodById(pool, id)
	if err != nil {
		return models.PaymentMethods{}, err
	}

	return pm, nil
}

func DeletePaymentMethod(pool *pgxpool.Pool, id int) error {
	result, err := pool.Exec(context.Background(),
		"DELETE FROM payment_methods WHERE id=$1", id,
	)

	if err != nil {
		return fmt.Errorf("error: failed to delete payment method: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("error: payment method not found")
	}

	return nil
}