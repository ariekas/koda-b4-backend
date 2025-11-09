package models

import (
	"time"
)

type User struct {
	Id         int        `json:"id"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Pic      *string    `json:"pic,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	Address    *string    `json:"address,omitempty"`
	Password   string     `json:"password,omitempty"`
	Role       string     `json:"role,omitempty"`
	Profileid  *int       `json:"profileid,omitempty"`
	Created_at *time.Time `json:"created_at,omitempty"`
	Updated_at *time.Time `json:"updated_at,omitempty"`
}

var InputNewRole struct{
	Role string `json:"role"`
}