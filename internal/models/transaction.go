package models

import (
	"time"
)

type TransactionItem struct {
	ProductID            int     `json:"productId" example:"1"`
	ProductName          string  `json:"productName" example:"Latte"`
	Quantity             int     `json:"quantity" example:"2"`
	Image                string  `json:"image"`
	PriceAtTime          float64 `json:"priceAtTime"`
	DiscountAtTime       float64 `json:"discountAtTime"`
	SizeID               int     `json:"sizeId,omitempty"`
	SizeName             string  `json:"sizeName,omitempty"`
	SizeAdditionalCost   float64 `json:"sizeAdditionalCost,omitempty"`
	VariantID            int     `json:"variantId,omitempty"`
	VariantName          string  `json:"variantName,omitempty"`
	VariantAdditionalCost float64 `json:"variantAdditionalCost,omitempty"`
	Subtotal             float64 `json:"subtotal" example:"50000"`
}

type Transaction struct {
	ID                 int               `json:"id"`
	InvoiceNum         string            `json:"invoiceNum"`
	UserID             int               `json:"userId"`
	UserFullname       string            `json:"userFullname" example:"Ari Eka Saputra"`
	UserAddress        string            `json:"userAddress" example:"Jl. Sudirman No. 10"`
	UserPhone          string            `json:"userPhone" example:"08123456789"`
	Status             string            `json:"status" example:"pending"`
	Subtotal           float64           `json:"subtotal" example:"85000"`
	TaxName            string            `json:"taxName,omitempty" example:"PPN"`
	TaxPercentage      float64           `json:"taxPercentage,omitempty" example:"10"`
	TaxAmount          float64           `json:"taxAmount" example:"8500"`
	DeliveryName       string            `json:"deliveryName" example:"Gojek"`
	DeliveryPrice      float64           `json:"deliveryPrice" example:"15000"`
	Total              float64           `json:"total" example:"108500"`
	PaymentMethod      string            `json:"paymentMethod" example:"Cash"`
	PaymentMethodImage *string           `json:"paymentMethodImage,omitempty"`
	Items              []TransactionItem `json:"items"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
}

type InputNewStatus struct {
	StatusTransactionID int `json:"statusTransactionsId" binding:"required"`
}

type PaginationResponseTransaction struct {
	Data       []Transaction     `json:"data"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"totalPages"`
	Links      map[string]string `json:"links"`
}

type TransactionInput struct {
	DeliveryID          int    `json:"deliveryId" binding:"required"`
	PaymentMethodID     int    `json:"paymentMethodId" binding:"required"`
	TaxID               int    `json:"taxId,omitempty"`
	StatusTransactionID int    `json:"statusTransactionsId,omitempty"`
	NameUser            string `json:"nameUser,omitempty"`
	AddressUser         string `json:"addressUser,omitempty"`
	PhoneUser           string `json:"phoneUser,omitempty"`
	EmailUser           string `json:"emailUser,omitempty"`
}

type CartItems struct {
	ProductID          int     `json:"productId"`
	ProductSizeID      int     `json:"productSizeId"`
	ProductVariantID   int     `json:"productVariantId"`
	Quantity           int     `json:"quantity"`
	ProductPrice       float64 `json:"productPrice"`
	DiscountPercentage float64 `json:"discountPercentage"`
	DiscountAmount     float64 `json:"discountAmount"`
	SizeCost           float64 `json:"sizeCost"`
	VariantCost        float64 `json:"variantCost"`
}

type PaymentMethod struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	ImagePayment *string   `json:"imagePayment"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type TransactionResponse struct {
	Invoice       string  `json:"invoice"`
	Total         float64 `json:"total"`
	PaymentStatus string  `json:"paymentStatus"`
}

type ProfileData struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type Status struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

type Tax struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}

type Delivery struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TransactionSummary struct {
	Subtotal      float64 `json:"subtotal"`
	TaxName       string  `json:"taxName,omitempty"`
	TaxPercentage float64 `json:"taxPercentage,omitempty"`
	TaxAmount     float64 `json:"taxAmount"`
	DeliveryName  string  `json:"deliveryName"`
	DeliveryPrice float64 `json:"deliveryPrice"`
	Total         float64 `json:"total"`
}