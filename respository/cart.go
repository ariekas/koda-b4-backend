package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddToCart(pool *pgxpool.Pool, userID int, input models.AddToCartInput) error {
	ctx := context.Background()
	var orderID int

	err := pool.QueryRow(ctx, `SELECT id FROM orders WHERE user_id=$1 AND status='pending'`, userID).Scan(&orderID)
	if err != nil {
		err := pool.QueryRow(ctx, `INSERT INTO orders (user_id, status, total) VALUES ($1,'pending',0) RETURNING id`, userID).Scan(&orderID)
		if err != nil {
			fmt.Println("Error: Create orders", err)
		}
	}

	var itemID int
	err = pool.QueryRow(ctx, `
		SELECT id FROM order_items 
		WHERE order_id=$1 AND product_id=$2 AND COALESCE(size_product_id,0)=$3 AND COALESCE(variant_id,0)=$4
	`, orderID, input.ProductID, input.SizeID, input.VariantID).Scan(&itemID)

	if err == nil {
		_, err := pool.Exec(ctx, `
			UPDATE order_items
			SET quantity = quantity + $1, subtotal = subtotal + $2, updated_at = NOW()
			WHERE id = $3
		`, input.Quantity, input.Subtotal, itemID)
		if err != nil {
			fmt.Println("Error: Update order items", err)
		}
	} else {
		_, err := pool.Exec(ctx, `
			INSERT INTO order_items (order_id, product_id, size_product_id, variant_id, quantity, subtotal)
			VALUES ($1,$2,$3,$4,$5,$6)
		`, orderID, input.ProductID, input.SizeID, input.VariantID, input.Quantity, input.Subtotal)
		if err != nil {
			fmt.Println("Error: Create order_items", err)
		}
	}

	_, err = pool.Exec(ctx, `
		UPDATE orders
		SET total = (SELECT COALESCE(SUM(subtotal),0) FROM order_items WHERE order_id=$1)
		WHERE id=$1
	`, orderID)
	if err != nil {
		fmt.Println("Error: Failed to getting order_items", err)
	}

	return nil
}

func GetUserCart(pool *pgxpool.Pool, userID int) (models.Cart, error) {
	ctx := context.Background()
	var cart models.Cart

	err := pool.QueryRow(ctx, `
		SELECT id, status, total FROM orders WHERE user_id=$1 AND status='pending'
	`, userID).Scan(&cart.OrderID, &cart.Status, &cart.Total)
	if err != nil {
		return cart, err
	}

	rows, err := pool.Query(ctx, `
		SELECT 
			oi.id, p.id, p.name, COALESCE(sp.id,0), COALESCE(sp.name,''), COALESCE(v.id,0), COALESCE(v.name,''), oi.quantity, oi.subtotal,
			COALESCE(pi.image,'') 
		FROM order_items oi
		JOIN products p ON oi.product_id=p.id
		LEFT JOIN size_product sp ON oi.size_product_id = sp.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		LEFT JOIN image_products pi ON pi.product_id = p.id
		WHERE oi.order_id=$1
		GROUP BY oi.id,p.id,p.name,sp.id,sp.name,v.id,v.name,pi.image
	`, cart.OrderID)
	if err != nil {
		return cart, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.ID, &item.ProductID, &item.ProductName, &item.SizeID, &item.SizeName, &item.VariantID, &item.VariantName, &item.Quantity, &item.Subtotal, &item.ImageURL)
		cart.Items = append(cart.Items, item)
	}
	return cart, nil
}

func Checkout(pool *pgxpool.Pool, userID int, paymentMethod string) error {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		UPDATE orders
		SET payment_method=$1, status='done', updated_at=NOW()
		WHERE user_id=$2 AND status='pending'
	`, paymentMethod, userID)

	return err
}

func GetUserCartProducts(pool *pgxpool.Pool, userID int) ([]models.CartItem, error) {
	var items []models.CartItem
	ctx := context.Background()

	rows, err := pool.Query(ctx, `
		SELECT 
			oi.id,
			p.id AS product_id,
			p.name AS product_name,
			COALESCE(sp.id,0) AS size_id,
			COALESCE(sp.name,'') AS size_name,
			COALESCE(v.id,0) AS variant_id,
			COALESCE(v.name,'') AS variant_name,
			oi.quantity,
			oi.subtotal,
			COALESCE(pi.image,'') AS image_url
		FROM order_items oi
		JOIN orders o ON o.id = oi.order_id
		JOIN products p ON oi.product_id = p.id
		LEFT JOIN size_product sp ON oi.size_product_id = sp.id
		LEFT JOIN variant v ON oi.variant_id = v.id
		LEFT JOIN image_products pi ON pi.product_id = p.id
		WHERE o.user_id = $1 AND o.status = 'pending'
		GROUP BY oi.id, p.id, p.name, sp.id, sp.name, v.id, v.name, pi.image
	`, userID)
	if err != nil {
		fmt.Println("Error: Failed to get data cart", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.SizeID,
			&item.SizeName,
			&item.VariantID,
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

