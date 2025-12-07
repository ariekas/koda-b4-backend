package models

import "time"

type Taxe struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Tax       float64   `json:"tax"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaxRequest struct {
	Name string  `json:"name" binding:"required"`
	Tax  float64 `json:"tax" binding:"required"`
}

type UpdateTaxRequest struct {
	Name string  `json:"name" binding:"required"`
	Tax  float64 `json:"tax" binding:"required"`
}
