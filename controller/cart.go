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

	err := respository.AddToCart(cc.Pool, userID,req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to add product to cart",
		})
		return
	}

	cart, err := respository.GetUserCart(cc.Pool, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch user cart",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Product added to cart successfully",
		Data:    cart,
	})
}

func (cc CartController) Checkout(ctx *gin.Context) {
	userID := middelware.GetUserFromToken(ctx)
	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	err := respository.Checkout(cc.Pool, userID, req.PaymentMethod)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Checkout failed",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Checkout successful",
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

	cartItems, err := respository.GetUserCartProducts(cc.Pool, userID)
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
