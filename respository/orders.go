package respository

import (
	"back-end-coffeShop/models"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PaginationResponseOrder struct {
	Data       []models.Order   `json:"data"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	Total      int              `json:"total"`
	TotalPages int              `json:"total_pages"`
	Links      map[string]string `json:"links"`
}

func GetOrders(pool *pgxpool.Pool, page int) (PaginationResponseOrder, error) {
	var orders []models.Order
	limit := 50
	offset := (page - 1) * limit

	var total int
	err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM orders").Scan(&total)
	if err != nil {
		fmt.Println("Error counting orders:", err)
	}

	rows, err := pool.Query(context.Background(), `
	SELECT
		o.id AS order_id,
		o.created_at,
		o.status,
		o.total,
		json_agg(
			json_build_object(
				'product_name', p.name
			)
		) FILTER (WHERE p.id IS NOT NULL) AS products
	FROM orders o
	LEFT JOIN orderitems oi ON o.id = oi.orderid
	LEFT JOIN product p ON p.id = oi.productid
	GROUP BY o.id, o.created_at, o.status, o.total
	ORDER BY o.id
	OFFSET $1 LIMIT $2
	`, offset, limit)
	if err != nil {
		fmt.Println("Error: Failed get data orders", err)
	}

	for rows.Next() {
		var o models.Order
		var productItems []byte

		err := rows.Scan(
			&o.ID,
			&o.CreatedAt,
			&o.Status,
			&o.Total,
			&productItems,
		)
		if err != nil {
			fmt.Println("Error scanning order:", err)
			continue
		}

		json.Unmarshal(productItems, &o.OrderItems)
		orders = append(orders, o)
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	links := make(map[string]string)

	if page > 1 {
		links["prev"] = fmt.Sprintf("/orders?page=%d", page-1)
	} else {
		links["prev"] = "null"
	}

	if page < totalPages {
		links["next"] = fmt.Sprintf("/orders?page=%d", page+1)
	} else {
		links["next"] = "null"
	}

	response := PaginationResponseOrder{
		Data:       orders,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}

	return response, nil
}

func GetOrderById(pool *pgxpool.Pool, orderId int) (models.Order, error) {
	var order models.Order

	query := `
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
	WHERE o.id = $1
	GROUP BY 
	  o.id, 
	  o.status, 
	  o.total,
	  u.fullname,
	  pr.address,
	  pr.phone,
	  o.paymentmethod,
	  d.type
	`

	row := pool.QueryRow(context.Background(), query, orderId)

	var orderItems []byte
	err := row.Scan(
		&order.ID,
		&order.Status,
		&order.Total,
		&order.UserFullname,
		&order.UserAddress,
		&order.UserPhone,
		&order.PaymentMethod,
		&order.DeliveryName,
		&orderItems,
	)
	if err != nil {
		fmt.Println("Error scanning order:", err)
		return models.Order{}, err
	}

	json.Unmarshal(orderItems, &order.OrderItems)

	return order, nil
}


func UpdateStatus(pool *pgxpool.Pool, orderId int, newStatus string) error {
	_, err := pool.Exec(context.Background(), "UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2", newStatus, orderId)

	if err != nil {
		fmt.Println("Error updating order status:", err)
	}

	return nil
}
