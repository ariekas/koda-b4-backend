package models

import "time"

type Discount struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Diskon float64 `json:"diskon"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}