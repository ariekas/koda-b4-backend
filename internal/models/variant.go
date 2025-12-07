package models

import "time"


type Variant struct {
	ID              int       `json:"id"`
	Name            string    `json:"name" `
	AdditionalCosts float64   `json:"additional_costs" `
	CreatedAt       time.Time `json:"created_at" `
	UpdatedAt       time.Time `json:"updated_at" `
}


type CreateVariant struct {
	Name            string  `json:"name"`
	AdditionalCosts float64 `json:"additional_costs" `
}

type UpdateVariant struct {
	Name            string  `json:"name"`
	AdditionalCosts float64 `json:"additional_costs"`
}