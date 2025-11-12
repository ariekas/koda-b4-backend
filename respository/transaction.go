package respository

import (
	"back-end-coffeShop/models"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetTransactions(pool *pgxpool.Pool, page int, limit int) (models.PaginationResponseTransaction, error) {
	var transactions []models.Transaction
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM transactions").Scan(&total)
	if err != nil {
		fmt.Println("Error counting transactions:", err)
	}

	query := `
	SELECT 
		t.id,
		t.user_id,
		u.fullname,
		pr.address,
		pr.phone,
		t.status,
		t.total,
		t.payment_method,
		d.type AS delivery_name,
		json_agg(
			json_build_object(
				'product_id', ti.product_id,
				'product_name', p.name,
				'quantity', ti.quantity,
				'subtotal', ti.subtotal
			)
		) FILTER (WHERE ti.id IS NOT NULL) AS items,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN transaction_items ti ON t.id = ti.transaction_id
	LEFT JOIN products p ON p.id = ti.product_id
	LEFT JOIN users u ON u.id = t.user_id
	LEFT JOIN profile pr ON pr.id = u.profile_id
	LEFT JOIN deliverys d ON d.transaction_id = t.id
	GROUP BY t.id, u.fullname, pr.address, pr.phone, d.type
	ORDER BY t.id DESC
	OFFSET $1 LIMIT $2
	`

	rows, err := pool.Query(context.Background(), query, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data transaction", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction
		var itemsJSON []byte

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.UserFullname,
			&t.UserAddress,
			&t.UserPhone,
			&t.Status,
			&t.Total,
			&t.PaymentMethod,
			&t.DeliveryName,
			&itemsJSON,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			fmt.Println("Error scanning transaction:", err)
			continue
		}

		json.Unmarshal(itemsJSON, &t.Items)
		transactions = append(transactions, t)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := make(map[string]string)

	if page > 1 {
		links["prev"] = fmt.Sprintf("/transactions?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/transactions?page=%d", page+1)
	} else {
		links["next"] = "null"
	}

	response := models.PaginationResponseTransaction{
		Data:       transactions,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}

	return response, nil
}

func GetTransactionById(pool *pgxpool.Pool, transactionId int) (models.Transaction, error) {
	var t models.Transaction
	query := `
	SELECT 
		t.id,
		t.user_id,
		u.fullname,
		pr.address,
		pr.phone,
		t.status,
		t.total,
		t.payment_method,
		d.type AS delivery_name,
		json_agg(
			json_build_object(
				'product_id', ti.product_id,
				'product_name', p.name,
				'quantity', ti.quantity,
				'subtotal', ti.subtotal
			)
		) FILTER (WHERE ti.id IS NOT NULL) AS items,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN transaction_items ti ON t.id = ti.transaction_id
	LEFT JOIN products p ON p.id = ti.product_id
	LEFT JOIN users u ON u.id = t.user_id
	LEFT JOIN profile pr ON pr.id = u.profile_id
	LEFT JOIN deliverys d ON d.transaction_id = t.id
	WHERE t.id = $1
	GROUP BY t.id, u.fullname, pr.address, pr.phone, d.type
	`

	row := pool.QueryRow(context.Background(), query, transactionId)
	var itemsJSON []byte
	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.UserFullname,
		&t.UserAddress,
		&t.UserPhone,
		&t.Status,
		&t.Total,
		&t.PaymentMethod,
		&t.DeliveryName,
		&itemsJSON,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		fmt.Println("Error scanning transaction:", err)
	}

	json.Unmarshal(itemsJSON, &t.Items)
	return t, nil
}

func UpdateTransactionStatus(pool *pgxpool.Pool, transactionId int, newStatus string) error {
	_, err := pool.Exec(context.Background(),
		"UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2",
		newStatus, transactionId,
	)
	if err != nil {
		fmt.Println("Error updating transaction status:", err)
	}
	return err
}