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
		t.users_id AS user_id,
		u.fullname,
		pr.address,
		pr.phone,
		st.status,
		t.total,
		pm.name AS payment_method,
		d.name AS delivery_name,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN users u ON u.id = t.users_id
	LEFT JOIN profile pr ON pr.id = u.profile_id
	LEFT JOIN deliverys d ON d.id = t.deliverys_id
	LEFT JOIN payment_methods pm ON pm.id = t.payment_methods_id
	LEFT JOIN status_transactions st ON st.id = t.status_transactions_id
	ORDER BY t.id DESC
	OFFSET $1 LIMIT $2
	`

	rows, err := pool.Query(context.Background(), query, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data transaction:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction

        // Scan tanpa items
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
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			fmt.Println("Error scanning transaction:", err)
			continue
		}

		transactions = append(transactions, t)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	links := map[string]string{
		"prev": "null",
		"next": "null",
	}

	if page > 1 {
		links["prev"] = fmt.Sprintf("/transactions?page=%d", page-1)
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("/transactions?page=%d", page+1)
	}

	return models.PaginationResponseTransaction{
		Data:       transactions,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}, nil
}


func GetTransactionById(pool *pgxpool.Pool, transactionId int) (models.Transaction, error) {
	var t models.Transaction

	query := `
	SELECT 
		t.id,
		t.users_id,
		u.fullname,
		pr.address,
		pr.phone,
		st.status,
		t.total,
		pm.name AS payment_method,
		d.name AS delivery_name,
		JSON_AGG(
			JSON_BUILD_OBJECT(
				'product_id', ti.products_id,
				'product_name', p.name,
				'quantity', ti.quantity,
				'subtotal', ti.subtotal
			)
		) FILTER (WHERE ti.id IS NOT NULL) AS items,
		t.created_at,
		t.updated_at
	FROM transactions t
	LEFT JOIN transaction_items ti ON t.id = ti.transactions_id
	LEFT JOIN products p ON p.id = ti.products_id
	LEFT JOIN users u ON u.id = t.users_id
	LEFT JOIN profile pr ON pr.id = u.profile_id
	LEFT JOIN deliverys d ON d.id = t.deliverys_id
	LEFT JOIN payment_methods pm ON pm.id = t.payment_methods_id
	LEFT JOIN status_transactions st ON st.id = t.status_transactions_id
	WHERE t.id = $1
	GROUP BY 
		t.id,
		u.fullname,
		pr.address,
		pr.phone,
		st.status,
		pm.name,
		d.name
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
		return t, err
	}

	if itemsJSON != nil {
		json.Unmarshal(itemsJSON, &t.Items)
	} else {
		t.Items = []models.TransactionItem{}
	}

	return t, nil
}




func UpdateTransactionStatus(pool *pgxpool.Pool, transactionId int, newStatusID int) error {

    cmdTag, err := pool.Exec(context.Background(),
        "UPDATE transactions SET status_transactions_id = $1, updated_at = NOW() WHERE id = $2",
        newStatusID, transactionId,
    )

    if err != nil {
        fmt.Println("Error updating transaction status:", err)
        return err
    }

    if cmdTag.RowsAffected() == 0 {
        return fmt.Errorf("transaction with id %d not found", transactionId)
    }

    return nil
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
	ctx := context.Background()
	var id int

	profileData, err := GetProfileByUser(pool, userID)
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil data profil: %v", err)
	}

	nameUser := input.NameUser
	if nameUser == "" {
		nameUser = profileData.Fullname
	}

	addressUser := input.AddressUser
	if addressUser == "" {
		addressUser = profileData.Address
	}

	phoneUser := input.PhoneUser
	if phoneUser == "" {
		phoneUser = profileData.Phone
	}

	emailUser := input.EmailUser
	if emailUser == "" {
		emailUser = profileData.Email
	}

	if nameUser == "" || addressUser == "" || emailUser == "" {
		return 0, fmt.Errorf("data user tidak lengkap")
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (
			users_id, 
			deliverys_id, 
			payment_methods_id,
			status_transactions_id,
			name_user, 
			address_user, 	
			phone_user, 
			email_user, 
			total, 
			invoice_num
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		userID,
		input.DeliveryID,
		input.PaymentMethodID,
		input.StatusTransactionID,
		nameUser,
		addressUser,
		phoneUser,
		emailUser,
		total,
		invoice,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("gagal membuat transaksi: %v", err)
	}

	return id, nil
}

func CreateTransactionItem(pool *pgxpool.Pool, tx pgx.Tx, transactionID int, item models.CartItems, subtotal float64) error {
	ctx := context.Background()

	var currentStock int
	err := tx.QueryRow(ctx, `
		SELECT stock FROM products WHERE id = $1
	`, item.ProductID).Scan(&currentStock)
	if err != nil {
		return fmt.Errorf("failed to get current stock for product ID %d: %v", item.ProductID, err)
	}

	if currentStock < item.Quantity {
		return fmt.Errorf("insufficient stock for product ID %d: available %d, requested %d",
			item.ProductID, currentStock, item.Quantity)
	}

	_, err = tx.Exec(ctx, `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
	`, item.Quantity, item.ProductID)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %v", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO transaction_items 
			(transactions_id, products_id, quantity, subtotal, variant_products_id, size_products_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, transactionID, item.ProductID, item.Quantity, subtotal, item.VariantProductID, item.SizeProductID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction item: %v", err)
	}

	var remainingStock int
	err = tx.QueryRow(ctx, `
		SELECT stock FROM products WHERE id = $1
	`, item.ProductID).Scan(&remainingStock)
	if err != nil {
		return fmt.Errorf("failed to verify remaining stock: %v", err)
	}
	if remainingStock < 0 {
		return fmt.Errorf("stock below zero after transaction for product ID %d", item.ProductID)
	}

	return nil
}

func ClearCart(pool *pgxpool.Pool, tx pgx.Tx, userID int) error {
	_, err := tx.Exec(context.Background(), `DELETE FROM carts WHERE users_id=$1`, userID)

	return err
}

func GetProfileByUser(pool *pgxpool.Pool, userID int) (models.ProfileData, error) {
	var data models.ProfileData

	err := pool.QueryRow(context.Background(), `
		SELECT 
			u.fullname,
			u.email,
			COALESCE(p.address, '') AS address,
			COALESCE(p.phone, '') AS phone
		FROM users u
		LEFT JOIN profile p ON u.profile_id = p.id
		WHERE u.id = $1
	`, userID).Scan(&data.Fullname, &data.Email, &data.Address, &data.Phone)

	return data, err
}
