package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DiscountController struct {
	Pool *pgxpool.Pool
}


func (dc DiscountController) GetAll(ctx *gin.Context) {
		response, err := respository.GetDiscount(dc.Pool)

		if err != nil {
			ctx.JSON(401, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}
	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success getting data category",
		Data:    response,
	})
}
