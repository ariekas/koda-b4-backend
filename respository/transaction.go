package respository

import (
	"back-end-coffeShop/models"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5"
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

func GetCartTransaction(pool *pgxpool.Pool, userID int) ([]models.CartItems, error) {
	query := `
	SELECT c.products_id, c.variant_products_id, c.size_products_id, c.quantity,
		p.price AS product_price,
		COALESCE(v.additional_costs, 0) AS variant_cost,
		COALESCE(s.additional_costs, 0) AS size_cost
		FROM carts c
		JOIN products p ON c.products_id = p.id
		LEFT JOIN variant_products v ON c.variant_products_id = v.id
		LEFT JOIN size_products s ON c.size_products_id = s.id
		WHERE c.users_id = $1
	`
	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		fmt.Println("Error: Failed to get cart")
	}

	defer rows.Close()

	var items []models.CartItems

	for rows.Next() {
		var item models.CartItems

		rows.Scan(&item.ProductID, &item.VariantProductID, &item.SizeProductID, &item.Quantity, &item.ProductPrice, &item.VariantCost, &item.SizeCost)
		items = append(items, item)
	}
	return items, nil
}

func GetDelivery(pool *pgxpool.Pool, deliveryID int) (float64, error) {
	var price float64

	err := pool.QueryRow(context.Background(), "SELECT price FROM deliverys WHERE id=$1", deliveryID).Scan(&price)

	return price, err
}

func CreateTransaction(pool *pgxpool.Pool, userID int, input models.TransactionInput, total float64, invoice string, tx pgx.Tx) (int, error) {
	var id int

	err := tx.QueryRow(context.Background(), `
		INSERT INTO transactions (users_id, deliverys_id, payment_methods_id, 
			name_user, address_user, phone_user, email_user, total, payment_status, invoice_num)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', $9)
		RETURNING id
	`, userID, input.DeliveryID, input.PaymentMethodID, input.NameUser,
		input.AddressUser, input.PhoneUser, input.EmailUser, total, invoice).Scan(&id)

	return id, err
}

func CreateTransactionItem(pool *pgxpool.Pool,tx pgx.Tx, transactionID int, item models.CartItems, subtotal float64) error{
	_, err := tx.Exec(context.Background(), `
	INSERT INTO transaction_items (transactions_id, products_id, quantity, subtotal, variant_products_id, size_products_id)
	VALUES ($1, $2, $3, $4, $5, $6)
`, transactionID, item.ProductID, item.Quantity, subtotal, item.VariantProductID, item.SizeProductID)

return err
}

func ClearCart(pool *pgxpool.Pool,tx pgx.Tx, userID int) error {
	_, err := tx.Exec(context.Background(), `DELETE FROM carts WHERE users_id=$1`, userID)

	return err
}