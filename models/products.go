package models

import (
	"time"
)

type Product struct {
	Id                 int       `json:"id"`
	Name               string    `json:"name"`
	Price              float64   `json:"price"`
	Description        string    `json:"description"`
	Stock              int       `json:"stock"`
	IsFlashSale        bool      `json:"is_flashsale"`
	IsFavoriteProduct  bool      `json:"is_favorite_product"`
	CategoryProductId  int       `json:"category_products_id"`
	Image              string    `json:"image"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type SizeProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	ProductId int       `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VariantProduct struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
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
type ProductDetail struct {
	Product  Product         `json:"product"`
	Images   []ImageProduct  `json:"images"`
	Sizes    []SizeProduct   `json:"sizes"`
	Variants []VariantProduct `json:"variants"`
}

type ProductInput struct {
	Name               string  `json:"name" binding:"required,min=3"`
	Price              float64 `json:"price" binding:"required,gt=0"`
	Description        string  `json:"description" binding:"required"`
	ProductSize        string  `json:"product_size" binding:"required,oneof=Small Medium Large"`
	Stock              int     `json:"stock" binding:"required,gte=0"`
	IsFlashSale        *bool   `json:"is_flashsale"`
	Temperature        string  `json:"temperature" binding:"required,oneof=Hot Cold"`
	CategoryProductId  int     `json:"category_products_id" binding:"required,gt=0"`
}
