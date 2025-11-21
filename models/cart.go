package models



type AddToCartInput struct {
	ProductID int `json:"productId" binding:"required"`
	SizeID    int `json:"sizeId,omitempty"`
	VariantID int `json:"variantId,omitempty"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type Cart struct {
	UserID       int         `json:"userId"`
	Items        []CartItem  `json:"items"`
	OrderTotal   float64     `json:"orderTotal"`
	Subtotal     float64     `json:"subtotal"`
	DeliveryCost float64     `json:"deliveryCost"`
	Tax          float64     `json:"tax"`
}

type CartItem struct {
	ID             int     `json:"id"`
	ProductID      int     `json:"productId"`
	ProductName    string  `json:"productName"`
	SizeID         int     `json:"sizeId,omitempty"`
	SizeName       string  `json:"sizeName,omitempty"`
	VariantID      int     `json:"variantId,omitempty"`
	VariantName    string  `json:"variantName,omitempty"`
	Quantity       int     `json:"quantity"`
	Price          float64 `json:"price"`
	DiscountPrice  float64 `json:"discountPrice"`
	FinalPrice     float64 `json:"finalPrice"`  
	OrderTotal     float64 `json:"orderTotal"`  
	ImageURL       string  `json:"imageUrl,omitempty"`
	IsFlashSale    bool    `json:"isFlashSale"`
}

