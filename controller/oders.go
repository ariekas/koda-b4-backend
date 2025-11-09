package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type OrdersController struct{
	Conn *pgx.Conn
}

func (oc OrdersController) GetOrders(ctx *gin.Context){
	order, err  := respository.GetOrders(oc.Conn)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to getting orders",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "SUccess getting data orders",
		Data: order,
	})
}