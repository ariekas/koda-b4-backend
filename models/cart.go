package models


type AddToCart struct {
	ProductID  int     `json:"product_id"`
	SizeID     int     `json:"size_id"`
	VariantID  int     `json:"variant_id"`
	Quantity   int     `json:"quantity"`
	Subtotal   float64 `json:"subtotal"`
}

type Cart struct {
	OrderID int               `json:"order_id"`
	Status  string            `json:"status"`
	Total   float64           `json:"total"`
	Items   []CartItems `json:"items"`
}

type CartItems struct {
	ID          int     `json:"id"`
	ProductName string  `json:"product_name"`
	SizeName    string  `json:"size_name"`
	VariantName string  `json:"variant_name"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
	ImageURL    string  `json:"image_url"`
}
