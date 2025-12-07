package models

type CategoryProduct struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required,min=3,max=50"`
}