package respository

import (
	"back-end-coffeShop/internal/models"
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
		return models.PaginationResponseTransaction{}, fmt.Errorf("error counting transactions: %v", err)
	}

	query := `
	SELECT 
		t.id,
		t.invoice_num,
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
		return models.PaginationResponseTransaction{}, fmt.Errorf("failed to get transactions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t models.Transaction

		err := rows.Scan(
			&t.ID,
			&t.InvoiceNum,
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
		t.name_user,
		t.address_user,
		t.phone_user,
		st.status,
		t.subtotal,
		t.tax_amount,
		t.total,
		t.invoice_num,
		pm.name AS payment_method,
		pm.image_payment,
		d.name AS delivery_name,
		d.price AS delivery_price,
		tx.name AS tax_name,
		tx.tax AS tax_percentage,

		JSON_AGG(
			JSON_BUILD_OBJECT(
				'productId', ti.products_id,
				'productName', p.name,
				'quantity', ti.quantity,
				'image', (
					SELECT pi.image
					FROM product_images pi
					WHERE pi.products_id = p.id AND pi.is_primary = true
					LIMIT 1
				),
				'priceAtTime', ti.price_at_time,
				'discountAtTime', ti.discount_at_time,
				'sizeId', COALESCE(ps.size_products_id, 0),
				'sizeName', COALESCE(sp.name, ''),
				'sizeAdditionalCost', COALESCE(sp.additional_costs, 0),
				'variantId', COALESCE(pv.variant_products_id, 0),
				'variantName', COALESCE(vp.name, ''),
				'variantAdditionalCost', COALESCE(vp.additional_costs, 0),
				'subtotal', (ti.price_at_time - ti.discount_at_time + 
					COALESCE(sp.additional_costs, 0) + 
					COALESCE(vp.additional_costs, 0)) * ti.quantity
			)
		) FILTER (WHERE ti.id IS NOT NULL) AS items,

		t.created_at,
		t.updated_at

	FROM transactions t
	LEFT JOIN transaction_items ti ON t.id = ti.transactions_id
	LEFT JOIN products p ON p.id = ti.products_id
	LEFT JOIN product_sizes ps ON ps.id = ti.product_sizes_id
	LEFT JOIN size_products sp ON sp.id = ps.size_products_id
	LEFT JOIN product_variants pv ON pv.id = ti.product_variants_id
	LEFT JOIN variant_products vp ON vp.id = pv.variant_products_id
	LEFT JOIN deliverys d ON d.id = t.deliverys_id
	LEFT JOIN payment_methods pm ON pm.id = t.payment_methods_id
	LEFT JOIN status_transactions st ON st.id = t.status_transactions_id
	LEFT JOIN taxes tx ON tx.id = t.taxes_id
	WHERE t.id = $1
	GROUP BY 
		t.id,
		st.status,
		pm.name,
		pm.image_payment,
		d.name,
		d.price,
		tx.name,
		tx.tax
	`

	row := pool.QueryRow(context.Background(), query, transactionId)

	var itemsJSON []byte
	var deliveryPrice, taxPercentage float64
	var taxName string

	err := row.Scan(
		&t.ID,
		&t.UserID,
		&t.UserFullname,
		&t.UserAddress,
		&t.UserPhone,
		&t.Status,
		&t.Subtotal,
		&t.TaxAmount,
		&t.Total,
		&t.InvoiceNum,
		&t.PaymentMethod,
		&t.PaymentMethodImage,
		&t.DeliveryName,
		&deliveryPrice,
		&taxName,
		&taxPercentage,
		&itemsJSON,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return t, fmt.Errorf("error scanning transaction: %v", err)
	}

	t.DeliveryPrice = deliveryPrice
	t.TaxName = taxName
	t.TaxPercentage = taxPercentage

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
		return fmt.Errorf("error updating transaction status: %v", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("transaction with id %d not found", transactionId)
	}

	return nil
}

func GetCartTransaction(pool *pgxpool.Pool, userID int) ([]models.CartItems, error) {
	query := `
	SELECT 
		c.products_id,
		COALESCE(c.product_variants_id, 0) AS product_variant_id,
		COALESCE(c.product_sizes_id, 0) AS product_size_id,
		c.quantity,
		p.price AS product_price,
		COALESCE(d.discount_percentage, 0) AS discount_percentage,
		COALESCE(vp.additional_costs, 0) AS variant_cost,
		COALESCE(sp.additional_costs, 0) AS size_cost
	FROM carts c
	JOIN products p ON c.products_id = p.id
	LEFT JOIN discounts d ON p.discounts_id = d.id
	LEFT JOIN product_variants pv ON c.product_variants_id = pv.id
	LEFT JOIN variant_products vp ON pv.variant_products_id = vp.id
	LEFT JOIN product_sizes ps ON c.product_sizes_id = ps.id
	LEFT JOIN size_products sp ON ps.size_products_id = sp.id
	WHERE c.users_id = $1
	`

	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}
	defer rows.Close()

	var items []models.CartItems

	for rows.Next() {
		var item models.CartItems

		err := rows.Scan(
			&item.ProductID,
			&item.ProductVariantID,
			&item.ProductSizeID,
			&item.Quantity,
			&item.ProductPrice,
			&item.DiscountPercentage,
			&item.VariantCost,
			&item.SizeCost,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %v", err)
		}

		if item.DiscountPercentage > 0 {
			item.DiscountAmount = item.ProductPrice * (item.DiscountPercentage / 100)
		}

		items = append(items, item)
	}

	return items, nil
}

func GetDelivery(pool *pgxpool.Pool, deliveryID int) (float64, error) {
	var price float64

	err := pool.QueryRow(context.Background(),
		"SELECT price FROM deliverys WHERE id=$1", deliveryID).Scan(&price)

	if err != nil {
		return 0, fmt.Errorf("delivery not found: %v", err)
	}

	return price, nil
}

func GetTax(pool *pgxpool.Pool, taxID int) (models.Tax, error) {
	var tax models.Tax

	err := pool.QueryRow(context.Background(),
		"SELECT id, name, tax FROM taxes WHERE id=$1", taxID).
		Scan(&tax.ID, &tax.Name, &tax.Percentage)

	if err != nil {
		return tax, fmt.Errorf("tax not found: %v", err)
	}

	return tax, nil
}

func CreateTransaction(pool *pgxpool.Pool, userID int, input models.TransactionInput, subtotal, taxAmount, total float64, invoice string, tx pgx.Tx) (int, error) {
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

	statusID := input.StatusTransactionID
	if statusID == 0 {
		statusID = 1
	}

	taxID := input.TaxID
	if taxID == 0 {
		taxID = 1
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO transactions (
			users_id, 
			deliverys_id, 
			payment_methods_id,
			status_transactions_id,
			taxes_id,
			name_user, 
			address_user, 	
			phone_user, 
			email_user,
			subtotal,
			tax_amount,
			total, 
			invoice_num
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`,
		userID,
		input.DeliveryID,
		input.PaymentMethodID,
		statusID,
		taxID,
		nameUser,
		addressUser,
		phoneUser,
		emailUser,
		subtotal,
		taxAmount,
		total,
		invoice,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("gagal membuat transaksi: %v", err)
	}

	return id, nil
}

func CreateTransactionItem(pool *pgxpool.Pool, tx pgx.Tx, transactionID int, item models.CartItems) error {
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
		SET stock = stock - $1, updated_at = NOW()
		WHERE id = $2
	`, item.Quantity, item.ProductID)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %v", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO transaction_items 
			(transactions_id, products_id, quantity, price_at_time, discount_at_time, 
			 product_sizes_id, product_variants_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		transactionID,
		item.ProductID,
		item.Quantity,
		item.ProductPrice,
		item.DiscountAmount,
		item.ProductSizeID,
		item.ProductVariantID,
	)
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
	if err != nil {
		return fmt.Errorf("failed to clear cart: %v", err)
	}
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

	if err != nil {
		return data, fmt.Errorf("profile not found: %v", err)
	}

	return data, nil
}

func GetPaymentMethod(pool *pgxpool.Pool) ([]models.PaymentMethod, error) {
	rows, err := pool.Query(context.Background(), `
		SELECT id, name, image_payment, created_at, updated_at
		FROM payment_methods
		ORDER BY id ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment methods: %v", err)
	}
	defer rows.Close()

	var payments []models.PaymentMethod

	for rows.Next() {
		var pm models.PaymentMethod
		err := rows.Scan(
			&pm.Id,
			&pm.Name,
			&pm.ImagePayment,
			&pm.CreatedAt,
			&pm.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment method: %v", err)
		}

		payments = append(payments, pm)
	}

	return payments, nil
}

func GetAllStatus(pool *pgxpool.Pool) ([]models.Status, error) {
	rows, err := pool.Query(context.Background(),
		"SELECT id, status FROM status_transactions ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch status: %v", err)
	}
	defer rows.Close()

	var status []models.Status

	for rows.Next() {
		var s models.Status
		if err := rows.Scan(&s.Id, &s.Status); err != nil {
			return nil, fmt.Errorf("failed to scan status: %v", err)
		}
		status = append(status, s)
	}

	return status, nil
}

func GetAllDeliveries(pool *pgxpool.Pool) ([]models.Delivery, error) {
	rows, err := pool.Query(context.Background(),
		"SELECT id, name, price, created_at, updated_at FROM deliverys ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch deliveries: %v", err)
	}
	defer rows.Close()

	var deliveries []models.Delivery

	for rows.Next() {
		var d models.Delivery
		if err := rows.Scan(&d.ID, &d.Name, &d.Price, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan delivery: %v", err)
		}
		deliveries = append(deliveries, d)
	}

	return deliveries, nil
}

func GetAllTaxes(pool *pgxpool.Pool) ([]models.Tax, error) {
	rows, err := pool.Query(context.Background(),
		"SELECT id, name, tax FROM taxes ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch taxes: %v", err)
	}
	defer rows.Close()

	var taxes []models.Tax

	for rows.Next() {
		var t models.Tax
		if err := rows.Scan(&t.ID, &t.Name, &t.Percentage); err != nil {
			return nil, fmt.Errorf("failed to scan tax: %v", err)
		}
		taxes = append(taxes, t)
	}

	return taxes, nil
}
