package models

import "time"

type Deliverys struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDeliveryRequest struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required"`
}

type UpdateDeliveryRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}