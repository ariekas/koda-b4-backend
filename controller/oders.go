package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdersController struct {
	Pool pgxpool.Pool
}

// GetOrders godoc
// @Summary Get all orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /orders [get]
func (oc OrdersController) GetOrders(ctx *gin.Context) {
	order, err := respository.GetOrders(&oc.Pool)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to getting orders",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "SUccess getting data orders",
		Data:    order,
	})
}

// UpdateStatus godoc
// @Summary Update status order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body models.InputNewStatus true "New Status Data"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /orders/status/{id} [patch]
func (oc OrdersController) UpdateStatus(ctx *gin.Context) {
	id := ctx.Param("id")

	// CONVERT STRING KE INT
	orderId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error : failed to converd type data")
	}

	var input models.InputNewStatus
	err = ctx.ShouldBindJSON(&input)

	if err != nil {
		fmt.Println("Error: ", err)
	}
	err = respository.UpdateStatus(&oc.Pool, orderId, input.Status)

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
