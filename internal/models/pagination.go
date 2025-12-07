package models

type PaginationResponse struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"totalPages"`
	Links      map[string]string `json:"links"`
	Data       any  `json:"data"`
}