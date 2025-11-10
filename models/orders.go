package models

import (
	"database/sql"
)
type OrderItems struct {
	ProductName string  `json:"product_name" example:"Latte"`
	Quantity    int     `json:"quantity" example:"2"`
	Subtotal    float64 `json:"subtotal" example:"50000"`
}

type Order struct {
	ID            int          `json:"id"`
	Status        string       `json:"status" example:"paid"`
	Total         float64      `json:"total" example:"100000"`
	UserFullname  string       `json:"user_fullname" example:"Ari Eka Saputra"`
	UserAddress   string       `json:"user_address" example:"Jl. Sudirman No. 10"`
	UserPhone     string       `json:"user_phone" example:"08123456789"`
	PaymentMethod string       `json:"payment_method" example:"Cash"`
	DeliveryName  string       `json:"delivery_name" example:"Gojek"`
	OrderItems    []OrderItems `json:"order_items"`
	CreatedAt     sql.NullString    `json:"created_at" example:"2025-11-10T10:00:00Z"`
	UpdatedAt     sql.NullString    `json:"updated_at" example:"2025-11-10T10:00:00Z"`
}

type InputNewStatus struct {
	Status string `json:"status" example:"completed" binding:"required,oneof=pending processing completed cancelled"`
}
