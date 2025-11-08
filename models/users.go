package models

import (
	"time"
)

type User struct{
	Id int `json:"id"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	Profileid int `json:"profile_id"`
	Token string
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
