package models

import "time"

type Size struct {
	Id              int       `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	AdditionalCosts float64   `json:"additional_costs" db:"additional_costs"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}


type SizeResponse struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	AdditionalCosts float64   `json:"additional_costs"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}