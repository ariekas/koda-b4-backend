package models

import (
	"time"
)
type TransactionItem struct {
	ProductID   int     `json:"productId" example:"1"`
	ProductName string  `json:"productName" example:"Latte"`
	Quantity    int     `json:"quantity" example:"2"`
	Subtotal    float64 `json:"subtotal" example:"50000"`
}


type Transaction struct {
	ID            int               `json:"id"`
	UserID        int               `json:"userId"`
	UserFullname  string            `json:"userFullname" example:"Ari Eka Saputra"`
	UserAddress   string            `json:"userAddress" example:"Jl. Sudirman No. 10"`
	UserPhone     string            `json:"userPhone" example:"08123456789"`
	Status        string            `json:"status" example:"pending"`
	Total         float64           `json:"total" example:"100000"`
	PaymentMethod string            `json:"paymentMethod" example:"Cash"`
	DeliveryName  string            `json:"deliveryName" example:"Gojek"`
	Items         []TransactionItem `json:"items"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type InputNewStatus struct {
	StatusTransactionID int `json:"statusTransactionsId" binding:"required"`
}

type PaginationResponseTransaction struct {
	Data       []Transaction       `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"totalPages"`
	Links      map[string]string   `json:"links"`
}

type TransactionInput struct {
	DeliveryID          int    `json:"deliveryId" binding:"required"`
	PaymentMethodID     int    `json:"paymentMethodId" binding:"required"`
	StatusTransactionID int    `json:"statusTransactionsId" binding:"required"`
	NameUser            string `json:"nameUser"`
	AddressUser         string `json:"addressUser"`
	PhoneUser           string `json:"phoneUser"`
	EmailUser           string `json:"emailUser"`
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
	PaymentStatus string  `json:"paymentStatus"`
}


type ProfileData struct {
	Fullname string
	Email    string
	Address  string
	Phone    string
}