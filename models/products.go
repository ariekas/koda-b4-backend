package models

import (
	"time"
)

type Product struct {
	Id                int              `json:"id"`
	DiscountsId       *int             `json:"discounts_id"`
	Name              string           `json:"name"`
	Price             float64          `json:"price"`
	PriceDiscount     float64          `json:"price_discount"`
	Description       string           `json:"description"`
	Stock             int              `json:"stock"`
	IsFlashSale       bool             `json:"is_flashsale"`
	IsFavoriteProduct bool             `json:"is_favorite_product"`
	CategoryProductId int              `json:"category_products_id"`
	Images            []ImageProduct   `json:"images"`
	Sizes             []SizeProduct    `json:"sizes"`
	Variants          []VariantProduct `json:"variants"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type SizeProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	ProductId int       `json:"product_id"`
	AdditionalCosts float64 `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	AdditionalCosts float64 `json:"price"`
	ProductId int       `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ImageProduct struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
