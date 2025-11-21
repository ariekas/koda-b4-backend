package models



type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	SizeID    int `json:"sizeId,omitempty"`
	VariantID int `json:"variantId,omitempty"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type Cart struct {
	UserID int        `json:"userId"`
	Items  []CartItem `json:"items"`
	Total  float64    `jsoCartItemsn:"total"`
}

type CartItem struct {
	ID                  int     `json:"id"`
	ProductID           int     `json:"productId"`
	ProductName         string  `json:"productName"`

	SizeID              int     `json:"sizeId"`
	SizeName            string  `json:"sizeName"`
	SizeAdditionalCost  float64 `json:"sizeAdditionalCost"`

	VariantID           int     `json:"variantId"`
	VariantName         string  `json:"variantName"`
	VariantAdditionalCost float64 `json:"variantAdditionalCost"`

	IsFlashSale         bool    `json:"isFlashSale"`
	Price               float64 `json:"price"`
	PriceDiscounts      float64 `json:"priceDiscounts"`

	Quantity            int     `json:"quantity"`

	Subtotal            float64 `json:"subtotal"`
	OrderTotal          float64 `json:"orderTotal"`

	ImageURL            string  `json:"imageUrl"`
}