package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"
	"strconv"

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

func (oc OrdersController) UpdateStatus(ctx *gin.Context){
	id := ctx.Param("id")

	// CONVERT STRING KE INT
	orderId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error : failed to converd type data")
	}

	err = ctx.ShouldBindJSON(&models.InputNewStatus)

	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = respository.UpdateStatus(oc.Conn, orderId , models.InputNewStatus.Status)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to update status order",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success update status order",
	})
}