package models

import "time"

type OrderItems struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type Order struct {
	ID        int            `json:"id"`
	CreatedAt *time.Time       `json:"created_at"`
	Status    string         `json:"status"`
	Total     float64        `json:"total"`
	OrderItems  []OrderItems `json:"orderItems"`
}

var InputNewStatus struct{
	Status string `json:"status"`
}