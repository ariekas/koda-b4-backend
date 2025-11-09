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
  o.status,
  o.total,
  u.fullname AS user_fullname,
  pr.address AS user_address,
  pr.phone AS user_phone,
  o.paymentmethod,
  d.type AS delivery_name,
  json_agg(
    json_build_object(
      'product_name', p.name,
      'quantity', oi.quantity,
      'subtotal', oi.subtotal
    )
  ) FILTER (WHERE oi.id IS NOT NULL) AS products
FROM orders o
LEFT JOIN orderitems oi ON o.id = oi.orderid
LEFT JOIN product p ON p.id = oi.productid
LEFT JOIN users u ON u.id = o.userid
LEFT JOIN profile pr ON pr.id = u.profileid
LEFT JOIN delivery d ON d.orderid = o.id
GROUP BY 
  o.id, 
  o.status, 
  o.total,
  u.fullname,
  pr.address,
  pr.phone,
  o.paymentmethod,
  d.type
ORDER BY o.id;
`)

	if err != nil {
		fmt.Println("Error: Failed get data orders", err)
	}

	for rows.Next() {
		var o models.Order
		var orderItems []byte
		err := rows.Scan(
			&o.ID,
			&o.Status,
			&o.Total,
			&o.UserFullname,
			&o.UserAddress,
			&o.UserPhone,
			&o.PaymentMethod,
			&o.DeliveryName,
			&orderItems,
		)
		if err != nil {
			fmt.Println("Error scanning order:", err)
		}

		json.Unmarshal(orderItems, &o.OrderItems)

		orders = append(orders, o)
	}

	return orders, nil
}

func UpdateStatus(conn *pgx.Conn, orderId int, newStatus string) error {
	_, err := conn.Exec(context.Background(), "UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2", newStatus, orderId)

	if err != nil {
		fmt.Println("Error updating order status:", err)
	}

	return nil
}
