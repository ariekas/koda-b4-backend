package models

import (
	"time"
)

type Product struct{
	Id int
	Name string
	Price float64
	Description string
	Productsize string
	Stock int
	Isflashsale *bool
	Tempelatur string
	Category_productid int
	Created_at time.Time
	Updated_at time.Time
}