package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ProductController struct{
	Conn *pgx.Conn
}

func (pc ProductController) GetProducts(ctx *gin.Context){
	produts, err := respository.GetProducts(pc.Conn)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Failed to getting data products",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Sucess getting data products",
		Data: produts,
	})
}

func (pc ProductController) CreateProduct(ctx *gin.Context){
	product := respository.Create(ctx, pc.Conn)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success Create product",
		Data: product,
	})
}

func (pc ProductController) EditProduct(ctx *gin.Context){
	newProduct, err := respository.Edit(pc.Conn, ctx)

	if err != nil {
		if err.Error() == "product not found" {
			ctx.JSON(404, models.Response{
				Success: false,
				Message: "Error : Product not found",
				Data:    nil,
			})
			return
		}

		ctx.JSON(500, models.Response{
			Success: false,
			Message: "Failed to edit product",
			Data:    nil,
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success edit product",
		Data: newProduct,
	})
}

func (pc ProductController) DeleteProduct(ctx *gin.Context){
	err := respository.Delete(pc.Conn, ctx)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Error: Failed to get product",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Sucess deleted product",
	})
}