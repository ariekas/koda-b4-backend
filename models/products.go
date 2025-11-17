package models

import (
	"time"
)

type Product struct {
	Id                int              `json:"id"`
	DiscountsId       *int             `json:"discountsId"`
	Name              string           `json:"name"`
	Price             float64          `json:"price"`
	PriceDiscount     float64          `json:"priceDiscount"`
	Description       string           `json:"description"`
	Stock             int              `json:"stock"`
	IsFlashSale       bool             `json:"isFlashsale"`
	IsFavoriteProduct bool             `json:"isFavoriteProduct"`
	CategoryProductId int              `json:"categoryProductsId"`
	Images            []ImageProduct   `json:"images"`
	Sizes             []SizeProduct    `json:"sizes"`
	Variants          []VariantProduct `json:"variants"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

type SizeProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	ProductId int       `json:"productId"`
	AdditionalCosts float64 `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type VariantProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	AdditionalCosts float64 `json:"price"`
	ProductId int       `json:"productId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ImageProduct struct {
	Id        int       `json:"id"`
	ProductId int       `json:"productId"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
