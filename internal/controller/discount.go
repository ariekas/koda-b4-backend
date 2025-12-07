package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"strconv"

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

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get discount data",
		Result:  response,
	})
}

func (dc DiscountController) Create(ctx *gin.Context) {
	var req models.Discount

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	response, err := respository.CreateDiscount(dc.Pool, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success create discount",
		Result:  response,
	})
}

func (dc DiscountController) Detail(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid discount ID",
		})
		return
	}

	response, err := respository.DetailDiscount(dc.Pool, id)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get discount detail",
		Result:  response,
	})
}

func (dc DiscountController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid discount ID",
		})
		return
	}

	var req models.Discount

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	response, err := respository.UpdateDiscount(dc.Pool, id, req)
	if err != nil {
		ctx.JSON(500, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success update discount",
		Result:  response,
	})
}

func (dc DiscountController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Invalid discount ID",
		})
		return
	}

	err = respository.DeleteDiscount(dc.Pool, id)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success delete discount",
	})
}
