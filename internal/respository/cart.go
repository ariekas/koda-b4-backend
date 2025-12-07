package respository

import (
	"back-end-coffeShop/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func AddToCart(pool *pgxpool.Pool, userId int, productId int, sizeId int, variantId int, quantity int) error {
	ctx := context.Background()

	var basePrice float64
	err := pool.QueryRow(ctx, `SELECT price FROM products WHERE id=$1`, productId).Scan(&basePrice)
	if err != nil {
		return fmt.Errorf("product not found: %v", err)
	}

	var sizeCost float64
	var productSizeId int
	if sizeId != 0 {
		err := pool.QueryRow(ctx, `
			SELECT ps.id, sp.additional_costs 
			FROM product_sizes ps
			JOIN size_products sp ON ps.size_products_id = sp.id
			WHERE ps.products_id = $1 AND ps.size_products_id = $2
		`, productId, sizeId).Scan(&productSizeId, &sizeCost)

		if err != nil {
			return fmt.Errorf("size not available for this product: %v", err)
		}
	}

	var variantCost float64
	var productVariantId int
	if variantId != 0 {
		err := pool.QueryRow(ctx, `
			SELECT pv.id, vp.additional_costs 
			FROM product_variants pv
			JOIN variant_products vp ON pv.variant_products_id = vp.id
			WHERE pv.products_id = $1 AND pv.variant_products_id = $2
		`, productId, variantId).Scan(&productVariantId, &variantCost)

		if err != nil {
			return fmt.Errorf("variant not available for this product: %v", err)
		}
	}

	var existingCartID int
	err = pool.QueryRow(ctx, `
		SELECT id FROM carts
		WHERE users_id=$1 AND products_id=$2 
		AND product_sizes_id=$3 AND product_variants_id=$4
	`, userId, productId, productSizeId, productVariantId).Scan(&existingCartID)

	if err == nil {
		_, err := pool.Exec(ctx, `
			UPDATE carts
			SET quantity = quantity + $1, updated_at = NOW()
			WHERE id = $2
		`, quantity, existingCartID)

		if err != nil {
			return fmt.Errorf("failed to update cart: %v", err)
		}
		return nil
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO carts (users_id, products_id, product_sizes_id, product_variants_id, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`, userId, productId, productSizeId, productVariantId, quantity)

	if err != nil {
		return fmt.Errorf("failed to insert to cart: %v", err)
	}

	return nil
}

func GetUserCart(pool *pgxpool.Pool, userId int) ([]models.CartItem, error) {
	ctx := context.Background()
	var items []models.CartItem

	rows, err := pool.Query(ctx, `
		SELECT 
			c.id,
			p.id AS product_id,
			p.name AS product_name,

			COALESCE(ps.id, 0) AS product_size_id,
			COALESCE(sp.id, 0) AS size_id,
			COALESCE(sp.name, '') AS size_name,
			COALESCE(sp.additional_costs, 0) AS size_additional_cost,

			COALESCE(pv.id, 0) AS product_variant_id,
			COALESCE(vp.id, 0) AS variant_id,
			COALESCE(vp.name, '') AS variant_name,
			COALESCE(vp.additional_costs, 0) AS variant_additional_cost,

			c.quantity,
			COALESCE(pi.image, '') AS image_url,

			p.is_flashsale,
			p.price,
			COALESCE(d.discount_percentage, 0) AS discount_percentage

		FROM carts c
		JOIN products p ON c.products_id = p.id
		LEFT JOIN product_sizes ps ON c.product_sizes_id = ps.id
		LEFT JOIN size_products sp ON ps.size_products_id = sp.id
		LEFT JOIN product_variants pv ON c.product_variants_id = pv.id
		LEFT JOIN variant_products vp ON pv.variant_products_id = vp.id
		LEFT JOIN discounts d ON p.discounts_id = d.id
		LEFT JOIN LATERAL (
			SELECT image
			FROM product_images
			WHERE products_id = p.id AND is_primary = true
			LIMIT 1
		) pi ON true
		WHERE c.users_id = $1
		ORDER BY c.created_at DESC
	`, userId)

	if err != nil {
		return nil, fmt.Errorf("failed to query cart items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductSizeID,
			&item.SizeID,
			&item.SizeName,
			&item.SizeAdditionalCost,
			&item.ProductVariantID,
			&item.VariantID,
			&item.VariantName,
			&item.VariantAdditionalCost,
			&item.Quantity,
			&item.ImageURL,
			&item.IsFlashSale,
			&item.Price,
			&item.DiscountPercentage,
		)

		if err != nil {
			fmt.Println("failed to scan cart item:", err)
			continue
		}

		if item.DiscountPercentage > 0 {
			discount := item.Price * (item.DiscountPercentage / 100)
			item.DiscountPrice = item.Price - discount
		} else {
			item.DiscountPrice = item.Price
		}

		fmt.Println("DEBUG:", item.Price, item.DiscountPercentage, item.DiscountPrice,
			item.SizeAdditionalCost, item.VariantAdditionalCost, item.Quantity)

		item.TotalPrice = (item.DiscountPrice + item.SizeAdditionalCost + item.VariantAdditionalCost) * float64(item.Quantity)

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return items, nil
}

func DeleteCart(pool *pgxpool.Pool, userId int, cartId int) error {
	ctx := context.Background()

	result, err := pool.Exec(ctx, `
		DELETE FROM carts
		WHERE id = $1 AND users_id = $2
	`, cartId, userId)

	if err != nil {
		return fmt.Errorf("failed to delete cart item: %v", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("cart item not found or unauthorized")
	}

	return nil
}

func CountCart(pool *pgxpool.Pool, userId int) (int, error) {
	ctx := context.Background()
	var count int

	err := pool.QueryRow(ctx, `
        SELECT COUNT(carts)
        FROM carts
        WHERE users_id = $1
    `, userId).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to count cart: %v", err)
	}

	return count, nil
}
