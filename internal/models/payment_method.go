package models

import "time"

type PaymentMethods struct {
	Id           int    `json:"id"`
	Name         string `json:"name" binding:"required"`
	ImagePayment string `json:"image_payment"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

type PaymentMethodRequest struct {
	Name string `form:"name" binding:"required"`
}