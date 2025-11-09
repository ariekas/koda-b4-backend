package models

import "time"

type OrderItems struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Subtotal    float64 `json:"subtotal"`
}

type Order struct {
    ID           int
    Status       string
    Total        float64
    UserFullname string
    UserAddress  string
    UserPhone    string
    PaymentMethod string
    DeliveryName string
    OrderItems   []OrderItems
	Created_at time.Time
	Updated_at time.Time
}


var InputNewStatus struct{
	Status string `json:"status"`
}