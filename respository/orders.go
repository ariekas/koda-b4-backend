package respository

import (
	"back-end-coffeShop/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetOrders(conn *pgx.Conn) ([]models.Order, error) {
	var orders []models.Order
	rows, err := conn.Query(context.Background(), `
	SELECT
  o.id AS order_id,
  o.created_at,
  o.status,
  o.total,
  COALESCE(
    json_agg(
      json_build_object(
        'product_name', p.name,
        'quantity', oi.quantity,
        'subtotal', oi.subtotal
      )
    ) FILTER (WHERE oi.id IS NOT NULL),
    '[]'
  ) AS products
FROM orders o
LEFT JOIN orderitems oi ON o.id = oi.orderid
LEFT JOIN product p ON p.id = oi.productid
GROUP BY o.id, o.created_at, o.status, o.total
ORDER BY o.id;

`)

	if err != nil {
		fmt.Println("Error: Failed get data orders")
	}

	for rows.Next(){
		var o models.Order
		var orderItems []byte
		err := rows.Scan(&o.ID, &o.CreatedAt, &o.Status, &o.Total, &orderItems)
		
		if err != nil {
			fmt.Println("Error scanning order:", err)
		}

		json.Unmarshal(orderItems, &o.OrderItems)

		orders = append(orders, o)
	}

	return orders, nil
}