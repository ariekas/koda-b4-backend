package models

import (
	"time"
)

type User struct {
	Id         int        `json:"id"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Role       string     `json:"role"`
	ProfileID  *int       `json:"profile_id,omitempty"`
	Pic        *string    `json:"pic,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	Address    *string    `json:"address,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
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

type UpdateProfileRequest struct {
	Pic     *string `json:"pic" example:"https://example.com/profile.jpg"`
	Phone   *string `json:"phone" example:"+628123456789"`
	Address *string `json:"address" example:"Jl. Sudirman No. 123"`
}
