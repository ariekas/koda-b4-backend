package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AuthRegister struct{
	Conn *pgx.Conn
}

func (ar AuthRegister) Register(ctx *gin.Context) {
	user  := respository.Register(ctx, ar.Conn)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success register",
		Data: user,
	})

	
}