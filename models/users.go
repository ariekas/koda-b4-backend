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

type InputNewRoleUser struct {
	Role string `json:"role" example:"admin" binding:"required,oneof=admin user"`
}

type RegisterRequest struct {
	Fullname string `json:"fullname" example:"John Doe" binding:"required,min=3,max=100"`
	Email    string `json:"email" example:"john@example.com" binding:"required,email"`
	Password string `json:"password" example:"123456" binding:"required,min=6"`
	Role     string `json:"role" example:"customer" binding:"required,oneof=admin user"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com" binding:"required,email"`
	Password string `json:"password" example:"123456" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

var OtpForget = make(map[string]struct {
	Code      string
	ExpiresAt time.Time
})