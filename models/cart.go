package models



type AddToCartInput struct {
	ProductID int `json:"product_id" binding:"required"`
	SizeID    int `json:"size_id,omitempty"`
	VariantID int `json:"variant_id,omitempty"`
	Quantity  int `json:"quantity" binding:"required,min=1"`
}

type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
	Total  float64    `jsoCartItemsn:"total"`
}

type CartItem struct {
	ID           int     `json:"id"`
	ProductID    int     `json:"product_id"`
	ProductName  string  `json:"product_name"`
	SizeID       int     `json:"size_id,omitempty"`
	SizeName     string  `json:"size_name,omitempty"`
	VariantID    int     `json:"variant_id,omitempty"`
	VariantName  string  `json:"variant_name,omitempty"`
	Quantity     int     `json:"quantity"`
	Subtotal     float64 `json:"subtotal"`      
	ImageURL     string  `json:"image_url,omitempty"`
}
