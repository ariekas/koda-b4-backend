package respository

import (
	"back-end-coffeShop/models"
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
	if sizeId != 0 {
		if err := pool.QueryRow(ctx, `SELECT additional_costs FROM size_products WHERE id=$1`, sizeId).
			Scan(&sizeCost); err != nil {
			return fmt.Errorf("size not found: %v", err)
		}
	}

	var variantCost float64
	if variantId != 0 {
		if err := pool.QueryRow(ctx, `SELECT additional_costs FROM variant_products WHERE id=$1`, variantId).
			Scan(&variantCost); err != nil {
			return fmt.Errorf("variant not found: %v", err)
		}
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
			return fmt.Errorf("failed to update cart: %v", err)
		}
		return nil
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO carts (users_id, products_id, size_products_id, variant_products_id, quantity)
		VALUES ($1, $2, $3, $4, $5)
	`, userId, productId, sizeId, variantId, quantity)

	if err != nil {
		return fmt.Errorf("failed to insert to cart: %v", err)
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
    COALESCE(sp.additional_costs, 0) AS size_additional_cost,
    
    COALESCE(vp.id, 0) AS variant_id,
    COALESCE(vp.name, '') AS variant_name,
    COALESCE(vp.additional_costs, 0) AS variant_additional_cost,

    p.is_flashsale,
    p.price,
    p.price_discounts,

    c.quantity,

    (
        (CASE 
            WHEN p.price_discounts > 0 THEN p.price_discounts
            ELSE p.price
         END)
         + COALESCE(sp.additional_costs, 0)
         + COALESCE(vp.additional_costs, 0)
    ) * c.quantity AS subtotal,

    (CASE 
    WHEN p.price_discounts > 0 THEN p.price_discounts
    ELSE p.price
 END) * c.quantity AS order_total,


    COALESCE(pi.image, '') AS image_url

FROM carts c
JOIN products p ON c.products_id = p.id
LEFT JOIN size_products sp ON c.size_products_id = sp.id
LEFT JOIN variant_products vp ON c.variant_products_id = vp.id
LEFT JOIN LATERAL (
    SELECT image
    FROM product_images
    WHERE products_id = p.id
    ORDER BY id ASC
    LIMIT 1
) pi ON TRUE
WHERE c.users_id = $1

GROUP BY 
    c.id,
    p.id,
    p.name,
    sp.id,
    sp.name,
    sp.additional_costs,
    vp.id,
    vp.name,
    vp.additional_costs,
    p.is_flashsale,
    p.price,
    p.price_discounts,
    c.quantity,
    pi.image
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
			&item.SizeAdditionalCost,

			&item.VariantID,
			&item.VariantName,
			&item.VariantAdditionalCost,

			&item.IsFlashSale,
			&item.Price,
			&item.PriceDiscounts,

			&item.Quantity,

			&item.Subtotal,
			&item.OrderTotal,

			&item.ImageURL,
		)

		if err != nil {
			fmt.Println("failed to scan cart item:", err)
			continue
		}

		items = append(items, item)
	}
	
	if rows.Err() != nil {
		fmt.Println("error iterating rows:", rows.Err())
	}

	return items, nil
}

func DeleteCart(pool *pgxpool.Pool, userId int, cartId int) error {
	_, err := pool.Exec(context.Background(), `
	DELETE FROM carts
	WHERE id = $1 AND users_id = $2
`, cartId, userId)

fmt.Println(err)

return err
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
