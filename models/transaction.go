package models

import (
	"time"
)
type TransactionItem struct {
	ProductID   int     `json:"product_id" example:"1"`
	ProductName string  `json:"product_name" example:"Latte"`
	Quantity    int     `json:"quantity" example:"2"`
	Subtotal    float64 `json:"subtotal" example:"50000"`
}


type Transaction struct {
	ID            int               `json:"id"`
	UserID        int               `json:"user_id"`
	UserFullname  string            `json:"user_fullname" example:"Ari Eka Saputra"`
	UserAddress   string            `json:"user_address" example:"Jl. Sudirman No. 10"`
	UserPhone     string            `json:"user_phone" example:"08123456789"`
	Status        string            `json:"status" example:"pending"`
	Total         float64           `json:"total" example:"100000"`
	PaymentMethod string            `json:"payment_method" example:"Cash"`
	DeliveryName  string            `json:"delivery_name" example:"Gojek"`
	Items         []TransactionItem `json:"items"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type InputNewStatus struct {
	Status string `json:"status" example:"completed" binding:"required,oneof=pending processing completed cancelled"`
}

type PaginationResponseTransaction struct {
	Data       []Transaction       `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
	Links      map[string]string   `json:"links"`
}