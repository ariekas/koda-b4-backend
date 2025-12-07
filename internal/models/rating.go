package models

import "time"

type Rating struct {
	ID            int       `json:"id"`
	UsersID       int       `json:"users_id"`
	ProductsID    int       `json:"products_id"`
	TransactionsID int      `json:"transactions_id"`
	Rating        int       `json:"rating"`
	Review        string    `json:"review"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UnratedProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}

type CreateRatingRequest struct {
	ProductsID     int    `json:"products_id" binding:"required"`
	TransactionsID int    `json:"transactions_id" binding:"required"`
	Rating         int    `json:"rating" binding:"required,min=1,max=5"`
	Review         string `json:"review"`
}

type UpdateRatingRequest struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Review string `json:"review"`
}
