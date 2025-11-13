package models

import "time"

type History struct{
	Id int `json:"id"`
	Invoice string  `json:"invoice"`
	Date time.Time `json:"date"`
	Status string `json:"status"`
	Total float64 `json:"total"`
	ImageProduct string `json:"image"`
}

type DetailHistory struct {
	Id int `json:"id"`
	Invoice string  `json:"invoice"`
	Date time.Time `json:"date"`
	Fullname string `json:"fullname"`
	Address string `json:"address"`
	Payment string `json:"payment"`
	Phone string `json:"phone"`
	Delivery string `json:"delivery"`
	Status string `json:"status"`
	Total float64 `json:"total"`
	HistoryItems []HistoryItems `json:"items"`
}

type HistoryItems struct{
	Id int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Price float64 `json:"price"`
	PriceDiscount float64 `json:"price discount"`
	Size string `json:"size"`
	Variant string `json:"variant"`
	Quantity int `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

type PaginationResponseHistory struct {
	Data       []History       `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
	Links      map[string]string   `json:"links"`
}