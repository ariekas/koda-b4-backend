package models



type AddToCartInput struct {
	ProductID int     `json:"product_id" binding:"required"`
	SizeID    int     `json:"size_id"`   
	VariantID int     `json:"variant_id"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	Subtotal  float64 `json:"subtotal" binding:"required"`
}

type Cart struct {
	OrderID int         `json:"order_id"`
	UserID  int         `json:"user_id,omitempty"`
	Status  string      `json:"status"`
	Total   float64     `json:"total"`
	Items   []CartItem  `json:"items"`
}

type CartItem struct {
	ID          int     `json:"id"`
	ProductID   int     `json:"product_id"`
	ProductName string  `json:"product_name"`
	SizeID      int     `json:"size_id,omitempty"`
	SizeName    string  `json:"size_name,omitempty"`
	VariantID   int     `json:"variant_id,omitempty"`
	VariantName string  `json:"variant_name,omitempty"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
	ImageURL    string  `json:"image_url,omitempty"`
}