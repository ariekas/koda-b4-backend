package respository

import (
	"back-end-coffeShop/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddToCart(pool *pgxpool.Pool, userId int, productId int, sizeId int, variantId int, quantity int) error {
	ctx := context.Background()

	var basePrice, sizeCost, variantCost float64

	err := pool.QueryRow(ctx, `
		SELECT price FROM products WHERE id=$1
	`, productId).Scan(&basePrice)
	if err != nil {
		 fmt.Println("product not found", err)
	}

	if sizeId != 0 {
		_ = pool.QueryRow(ctx, `
			SELECT additional_costs FROM size_products WHERE id=$1
		`, sizeId).Scan(&sizeCost)
	}

	if variantId != 0 {
		_ = pool.QueryRow(ctx, `
			SELECT additional_costs FROM variant_products WHERE id=$1
		`, variantId).Scan(&variantCost)
	}

	var existingCartID int
	err = pool.QueryRow(ctx, `
		SELECT id FROM carts
		WHERE users_id=$1 AND products_id=$2 AND size_products_id=$3 AND variant_products_id=$4
	`, userId, productId, sizeId, variantId).Scan(&existingCartID)

	if err == nil {
		_, err := pool.Exec(ctx, `
			UPDATE carts
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2
		`, quantity, existingCartID)
		if err != nil {
			 fmt.Println("failed to update cart", err)
		}
	} else {
		_, err := pool.Exec(ctx, `
			INSERT INTO carts (users_id, products_id, size_products_id, variant_products_id, quantity)
			VALUES ($1, $2, $3, $4, $5)
		`, userId, productId, sizeId, variantId, quantity)
		if err != nil {
			 fmt.Println("failed to add to cart", err)
		}
	}

	return nil
}

func GetUserCartProduct(pool *pgxpool.Pool, userId int) ([]models.CartItem, error) {
	ctx := context.Background()
	var items []models.CartItem

	rows, err := pool.Query(ctx, `
		SELECT 
			c.id,
			p.id AS product_id,
			p.name AS product_name,
			COALESCE(sp.id, 0) AS size_id,
			COALESCE(sp.name, '') AS size_name,
			COALESCE(vp.id, 0) AS variant_id,
			COALESCE(vp.name, '') AS variant_name,
			c.quantity,
			(p.price 
				+ COALESCE(sp.additional_costs, 0) 
				+ COALESCE(vp.additional_costs, 0)
			) * c.quantity AS subtotal,
			COALESCE(pi.image, '') AS image_url
		FROM carts c
		JOIN products p ON c.products_id = p.id
		LEFT JOIN size_products sp ON c.size_products_id = sp.id
		LEFT JOIN variant_products vp ON c.variant_products_id = vp.id
		LEFT JOIN product_images pi ON pi.products_id = p.id
		WHERE c.users_id = $1
		GROUP BY 
			c.id, p.id, p.name, sp.id, sp.name, vp.id, vp.name, pi.image, 
			p.price, sp.additional_costs, vp.additional_costs, c.quantity
	`, userId)

	if err != nil {
		fmt.Println("failed to query cart items", err)
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
			fmt.Println("failed to scan cart item:", err)
		}
		items = append(items, item)
	}
	if rows.Err() != nil {
		fmt.Println("error iterating rows:", rows.Err())
	}

	return items, nil
}
