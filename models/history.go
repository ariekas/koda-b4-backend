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

type PaginationResponseHistory struct {
	Data       []History       `json:"data"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	Total      int                 `json:"total"`
	TotalPages int                 `json:"total_pages"`
	Links      map[string]string   `json:"links"`
}