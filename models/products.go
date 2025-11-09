package models

import (
	"time"
)

type Product struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	Description string `json:"description"`
	Productsize string `json:"productSize"`
	Stock int `json:"stock"`
	Isflashsale *bool `json:"isFlashsale"`
	Tempelatur string `json:"tempelatur"`
	Category_productid int `json:"category_productId"`
	Created_at time.Time `json:"creaed_at"`
	Updated_at time.Time `json:"updated_at"`
}

type ImageProduct struct {
	Id         int       `json:"id"`
	Productid  int       `json:"productId"`
	Image      string    `json:"image"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}