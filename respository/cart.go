package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddToCart(pool *pgxpool.Pool, userId int, productId int, sizeId int, variantId int, quantity int, subtotal float64) error {
	var orderId int
	err := pool.QueryRow(context.Background(), `SELECT id FROM orders WHERE user_id=$1 AND status='pending'`, userId).Scan(&orderId)
	if err != nil {
		err := pool.QueryRow(context.Background(), `
			INSERT INTO orders (user_id, status, total) VALUES ($1, 'pending', 0) RETURNING id
		`, userId).Scan(&orderId)
		if err != nil {
			fmt.Println("Error: Create orders", err)
		}
	}

	var existingItemId int
	err = pool.QueryRow(context.Background(), `
		SELECT id FROM order_items 
		WHERE order_id=$1 AND product_id=$2 AND size_product_id=$3 AND variant_id=$4
	`, orderId, productId, sizeId, variantId).Scan(&existingItemId)

	if err == nil {
		_, err := pool.Exec(context.Background(), `
			UPDATE order_items 
			SET quantity = quantity + $1, subtotal = subtotal + $2, updated_at = NOW()
			WHERE id = $3
		`, quantity, subtotal, existingItemId)
		if err != nil {
			fmt.Println("Error: Update order items", err)
		}
	} else {
		_, err := pool.Exec(context.Background(), `
			INSERT INTO order_items (order_id, product_id, size_product_id, variant_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, orderId, productId, sizeId, variantId, quantity, subtotal)
		if err != nil {
			fmt.Println("Error: Create order_items", err)
		}
	}

	_, err = pool.Exec(context.Background(), `
		UPDATE orders 
		SET total = (SELECT COALESCE(SUM(subtotal),0) FROM order_items WHERE order_id=$1)
		WHERE id=$1
	`, orderId)

	if err != nil {
		fmt.Println("Error: Failed to getting order_items", err)
	}

	return  nil
}

func GetUserCart(pool *pgxpool.Pool, userId int) (models.Cart, error) {
	ctx := context.Background()
	var cart models.Cart

	err := pool.QueryRow(ctx, `
		SELECT id, total, status
		FROM orders 
		WHERE user_id=$1 AND status='pending'
	`, userId).Scan(&cart.OrderID, &cart.Total, &cart.Status)
	if err != nil {
		return cart, err
	}

	rows, err := pool.Query(ctx, `
		SELECT 
			oi.id,
			p.name AS product_name,
			sp.name AS size_name,
			v.name AS variant_name,
			oi.quantity,
			oi.subtotal
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN size_product sp ON oi.size_product_id = sp.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		WHERE oi.order_id=$1
	`, cart.OrderID)
	if err != nil {
		return cart, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItems
		rows.Scan(&item.ID, &item.ProductName, &item.SizeName, &item.VariantName, &item.Quantity, &item.Subtotal)
		cart.Items = append(cart.Items, item)
	}
	return cart, nil
}

func Checkout(pool *pgxpool.Pool, userId int, paymentMethod string) error {
	ctx := context.Background()

	_, err := pool.Exec(ctx, `
		UPDATE orders
		SET payment_method=$1, status='done', updated_at=NOW()
		WHERE user_id=$2 AND status='pending'
	`, paymentMethod, userId)

	return err
}