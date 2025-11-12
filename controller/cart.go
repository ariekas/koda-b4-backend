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

func (cc CartController) AddCart(ctx *gin.Context) {
	userId := middelware.GetUserFromToken(ctx)

	var req models.AddToCart
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Error: Invalid request body",
		})
		return
	}

	item, status, orderId, err := respository.AddToCart(
		cc.Pool,
		userId,
		req.ProductID,
		req.SizeID,
		req.VariantID,
		req.Quantity,
	)

	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Error: Failed to create cart",
		})
		return
	}

	response := gin.H{
		"order_id":     orderId,
		"status":       status,
		"product_name": item.ProductName,
		"variant_name": item.VariantName,
		"size_name":    item.SizeName,
		"quantity":     item.Quantity,
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success to create cart",
		Data:    response,
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

func (cc CartController) GetCart(ctx *gin.Context) {
	userId := middelware.GetUserFromToken(ctx)

	if userId == 0 {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Unauthorized: Invalid token",
		})
		return
	}

	Cartitems, err := respository.GetUserCartProduct(cc.Pool, userId)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Error: Failed to get user cart",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get cart",
		Data:    Cartitems,
	})
}
