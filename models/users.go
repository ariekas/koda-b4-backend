package models

import (
	"database/sql"
)

type User struct{
	Id int `json:"id"`
	Fullname string `json:"fullname"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	Profileid int `json:"profile_id"`
	Created_at sql.NullTime `json:"created_at"`
	Updated_at sql.NullTime `json:"updated_at"`
}
