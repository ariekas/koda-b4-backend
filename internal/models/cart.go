package models

type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	SizeID    int `json:"sizeId,omitempty"`
	VariantID int `json:"variantId,omitempty"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type CartItem struct {
	ID                    int     `json:"id"`
	ProductID             int     `json:"productId"`
	ProductName           string  `json:"productName"`
	ProductSizeID         int     `json:"productSizeId,omitempty"`
	SizeID                int     `json:"sizeId,omitempty"`
	SizeName              string  `json:"sizeName,omitempty"`
	SizeAdditionalCost    float64 `json:"sizeAdditionalCost,omitempty"`
	ProductVariantID      int     `json:"productVariantId,omitempty"`
	VariantID             int     `json:"variantId,omitempty"`
	VariantName           string  `json:"variantName,omitempty"`
	VariantAdditionalCost float64 `json:"variantAdditionalCost,omitempty"`
	Quantity              int     `json:"quantity"`
	ImageURL              string  `json:"imageUrl,omitempty"`
	IsFlashSale           bool    `json:"isFlashsale"`
	Price                 float64 `json:"price"`
	DiscountPercentage    float64 `json:"discountPercentage,omitempty"`
	DiscountPrice         float64 `json:"discountPrice"`
	TotalPrice            float64 `json:"totalPrice"`
}
