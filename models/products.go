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

type ProductInput struct {
	Name               string  `json:"name" binding:"required,min=3" example:"Cappuccino"`
	Price              float64 `json:"price" binding:"required,gt=0" example:"25000"`
	Description        string  `json:"description" binding:"required" example:"Coffee with milk and foam"`
	Productsize        string  `json:"productSize" binding:"required,oneof=Small Medium Large" example:"Medium"`
	Stock              int     `json:"stock" binding:"required,gte=0" example:"20"`
	Isflashsale        *bool   `json:"isFlashsale" example:"false"`
	Tempelatur         string  `json:"tempelatur" binding:"required,oneof=Hot Cold" example:"Hot"`
	Category_productid int     `json:"category_productId" binding:"required,gt=0" example:"1"`
}

type ImageProduct struct {
	Id         int       `json:"id" binding:"required"`
	Productid  int       `json:"productId" binding:"required"`
	Image      string    `json:"image" binding:"required"`
	Created_at time.Time `json:"created_at" binding:"required"`
	Updated_at time.Time `json:"updated_at" binding:"required"`
}