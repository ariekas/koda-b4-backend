package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryProductController struct {
	Pool *pgxpool.Pool
}

func (cpc CategoryProductController) GetAll(ctx *gin.Context){
	categorys, err := respository.GetCategorys(cpc.Pool)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to getting categorys",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success getting data category",
		Data: categorys,
	})
}

func (cpc CategoryProductController) Create(ctx *gin.Context){
	categorys, err := respository.CreateCategory(cpc.Pool, ctx)
	if err != nil{
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Failed to create category",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success create category",
		Data: categorys,
	})
}

func (cpc CategoryProductController) GetByID(ctx *gin.Context){
	category, err := respository.GetCategoryById(cpc.Pool, ctx)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error get category",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success get data category",
		Data: category,
	})
}

func (cpc CategoryProductController) Edit(ctx *gin.Context){
	newCategory, err := respository.EditCategory(cpc.Pool, ctx)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error edit category",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success edit category",
		Data: newCategory,
	})
}

func (cpc CategoryProductController) Delete(ctx *gin.Context){
	err := respository.DeleteCategory(cpc.Pool, ctx)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error delete category",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success deleting category",
	})
}