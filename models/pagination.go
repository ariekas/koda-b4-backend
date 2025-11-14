package models

type PaginationResponse struct {
	Data       any  `json:"data"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
	Links      map[string]string `json:"links"`
}