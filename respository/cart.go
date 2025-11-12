package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddToCart(pool *pgxpool.Pool, userId int, productId int, sizeId int, variantId int, quantity int) (models.CartItems, string, int, error) {
	ctx := context.Background()
	var orderId int
	var status string = "pending"

	err := pool.QueryRow(ctx, `
		SELECT id FROM orders WHERE user_id=$1 AND status='pending'
	`, userId).Scan(&orderId)
	if err != nil {
		err = pool.QueryRow(ctx, `
			INSERT INTO orders (user_id, status, total)
			VALUES ($1, 'pending', 0)
			RETURNING id
		`, userId).Scan(&orderId)
		if err != nil {
			fmt.Println("Error: Create orders", err)
			return models.CartItems{}, status, 0, err
		}
	}

	var existingItemId int
	err = pool.QueryRow(ctx, `
		SELECT id FROM order_items
		WHERE order_id=$1 AND product_id=$2 AND size_product_id=$3 AND variant_id=$4
	`, orderId, productId, sizeId, variantId).Scan(&existingItemId)

	if err == nil {
		_, err = pool.Exec(ctx, `
			UPDATE order_items
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2
		`, quantity, existingItemId)
		if err != nil {
			fmt.Println("Error: Update order_items", err)
			return models.CartItems{}, status, orderId, err
		}
	} else {
		_, err = pool.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, size_product_id, variant_id, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`, orderId, productId, sizeId, variantId, quantity)
		if err != nil {
			fmt.Println("Error: Create order_items", err)
			return models.CartItems{}, status, orderId, err
		}
	}

	var item models.CartItems
	err = pool.QueryRow(ctx, `
		SELECT 
			p.name AS product_name,
			COALESCE(v.name, '') AS variant_name,
			COALESCE(sp.name, '') AS size_name,
			oi.quantity
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		LEFT JOIN size_product sp ON oi.size_product_id = sp.id
		WHERE oi.order_id=$1 AND oi.product_id=$2 AND oi.size_product_id=$3 AND oi.variant_id=$4
		ORDER BY oi.id DESC
		LIMIT 1
	`, orderId, productId, sizeId, variantId).Scan(
		&item.ProductName,
		&item.VariantName,
		&item.SizeName,
		&item.Quantity,
	)
	if err != nil {
		fmt.Println("Error: Failed to get last added item", err)
	}

	return item, status, orderId, nil
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

func GetUserCartProduct(pool *pgxpool.Pool, userId int) ([]models.CartItems, error) {
	var items []models.CartItems

	rows, err := pool.Query(context.Background(), `
		SELECT 
			oi.id,
			p.name AS product_name,
			COALESCE(sp.name, '') AS size_name,
			COALESCE(v.name, '') AS variant_name,
			oi.quantity,
			(oi.quantity * (
				p.price 
				+ COALESCE(sp.additional_costs, 0) 
				+ COALESCE(v.additional_costs, 0)
			)) AS subtotal,
			COALESCE(pi.image, '') AS image
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN size_product sp ON oi.size_product_id = sp.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		LEFT JOIN image_products pi ON pi.product_id = p.id
		WHERE o.user_id=$1 AND o.status='pending'
		GROUP BY oi.id, p.name, sp.name, v.name, pi.image, p.price, sp.additional_costs, v.additional_costs
	`, userId)
	if err != nil {
		fmt.Println("Error: Failed to get data cart", err)
	}

	for rows.Next() {
		var item models.CartItems
		err := rows.Scan(
			&item.ID,
			&item.ProductName,
			&item.SizeName,
			&item.VariantName,
			&item.Quantity,
			&item.Subtotal,
			&item.ImageURL,
		)
		if err != nil {
			fmt.Println("Error: Failed to scan cart items", err)
		}
		items = append(items, item)
	}

	return items, nil
}

