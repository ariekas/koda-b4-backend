package controller

import (
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartController struct {
	Pool *pgxpool.Pool
}

func (cc CartController) AddCart(ctx *gin.Context) {
	userID := middelware.GetUserFromToken(ctx)
	var req models.AddToCartInput

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	if err := respository.AddToCart(cc.Pool, userID, req.ProductID, req.SizeID, req.SizeID, req.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to add product to cart",
		})
		return
	}

	cartItems, err := respository.GetUserCartProduct(cc.Pool, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch updated cart",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Product added to cart successfully",
		Data:    cartItems,
	})
}

func (cc CartController) GetCart(ctx *gin.Context) {
	userID := middelware.GetUserFromToken(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "Unauthorized: Invalid token",
		})
		return
	}

	cartItems, err := respository.GetUserCartProduct(cc.Pool, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to get cart items",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Successfully fetched cart items",
		Data:    cartItems,
	})
}