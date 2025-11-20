package models

import (
	"mime/multipart"
	"time"
)

type User struct {
	Id         int        `json:"id"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Password   string     `json:"-"`
	Role       string     `json:"role"`
	ProfileID  *int       `json:"profileId,omitempty"`
	Pic        *string    `json:"pic,omitempty"`
	Phone      *string    `json:"phone,omitempty"`
	Address    *string    `json:"address,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
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
}

type LoginRequest struct {
	Email    string `json:"email" example:"john@example.com" binding:"required,email"`
	Password string `json:"password" example:"123456" binding:"required,min=6"`
}


var OtpForget = map[string]struct{
    Code       string
    ExpiresAt  time.Time
    Verified   bool
}{}

type UpdateProfileRequest struct {
	Pic     *string          `json:"pic"`
	Phone   *string          `json:"phone"`
	Address *string          `json:"address"`

	PicFile *multipart.FileHeader `json:"-"`
}

