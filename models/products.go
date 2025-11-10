package models

import (
	"database/sql"
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
	Created_at sql.NullString `json:"creaed_at"`
	Updated_at sql.NullString `json:"updated_at"`
}

type ProductInput struct {
	Name               string  `json:"name" example:"Cappuccino" binding:"required"`
	Price              float64 `json:"price" example:"25000" binding:"required"`
	Description        string  `json:"description" example:"Coffee with milk and foam"`
	Productsize        string  `json:"productSize" example:"Medium"`
	Stock              int     `json:"stock" example:"20"`
	Isflashsale        *bool   `json:"isFlashsale" example:"false"`
	Tempelatur         string  `json:"tempelatur" example:"Hot"`
	Category_productid int     `json:"category_productId" example:"1"`
}

type ImageProduct struct {
	Id         int       `json:"id"`
	Productid  int       `json:"productId"`
	Image      string    `json:"image"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}