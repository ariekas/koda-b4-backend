package controller

import (
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartController struct {
	Pool *pgxpool.Pool
}

func (cc CartController) AddCart(ctx *gin.Context)  {
	userId := middelware.GetUserFromToken(ctx)

	var req models.AddToCart

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed type request",
		})
		return
	}

	err := respository.AddToCart(cc.Pool, userId, req.ProductID, req.SizeID, req.VariantID, req.Quantity, req.Subtotal)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to create cart",
		})
		return
	}

	cart, err := respository.GetUserCart(cc.Pool, userId)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Error : Failed to get user",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success to create cart",
		Data: cart,
	})
}

func (cc CartController) Checkout(ctx *gin.Context) {
	userId := middelware.GetUserFromToken(ctx)
	var req struct {
		PaymentMethod string `json:"payment_method"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed type request",
		})
		return
	}

	err := respository.Checkout(cc.Pool, userId, req.PaymentMethod)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to checkout",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Checkout success",
	})
}