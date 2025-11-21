package models



type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	SizeID    int `json:"sizeId,omitempty"`
	VariantID int `json:"variantId,omitempty"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type CartItem struct {
	ID           int     `json:"id"`
	ProductID    int     `json:"productId"`
	ProductName  string  `json:"productName"`
	SizeID       int     `json:"sizeId,omitempty"`
	SizeName     string  `json:"sizeName,omitempty"`
	VariantID    int     `json:"variantId,omitempty"`
	VariantName  string  `json:"variantName,omitempty"`
	Quantity     int     `json:"quantity"`
	ImageURL     string  `json:"imageUrl,omitempty"`
	IsFlashSale bool `json:"is_flashsale"`
	Price float64 `json:"price"`
	DiscoutPrice float64 `json:"price_discounts"`
}