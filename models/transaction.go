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

type TransactionInput struct {
	DeliveryID      int     `json:"delivery_id" binding:"required"`
	PaymentMethodID int     `json:"payment_method_id" binding:"required"`
	NameUser        string  `json:"name_user"`
	AddressUser     string  `json:"address_user"`
	PhoneUser       string  `json:"phone_user"`
	EmailUser       string  `json:"email_user"`
}

type CartItems struct {
	ProductID        int
	VariantProductID int
	SizeProductID    int
	Quantity         int
	ProductPrice     float64
	VariantCost      float64
	SizeCost         float64
}


type TransactionResponse struct {
	Invoice       string  `json:"invoice"`
	Total         float64 `json:"total"`
	PaymentStatus string  `json:"payment_status"`
}


type ProfileData struct {
	Fullname string
	Email    string
	Address  string
	Phone    string
}