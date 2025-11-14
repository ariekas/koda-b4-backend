package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryProductController struct {
	Pool *pgxpool.Pool
}

func (cpc CategoryProductController) GetAll(ctx *gin.Context){
	categorys, err := respository.GetCategories(cpc.Pool)

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
		Data: categorys,
	})
}

func (cpc CategoryProductController) Create(ctx *gin.Context){
	
	var input models.CategoryProduct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	categorys, err := respository.CreateCategory(cpc.Pool, input)
	
	if input.Name == categorys.Name {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "error: category already exists",
		})
		return
	}

	if err != nil{
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
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
	id := ctx.Param("id")
	categoryId,_ := strconv.Atoi(id)
	category, err := respository.GetCategoryById(cpc.Pool, categoryId)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success get data category",
		Data: category,
	})
}

func (cpc CategoryProductController) Edit(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	var input models.CategoryProduct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	updatedCategory, err := respository.EditCategory(cpc.Pool, id, input)
	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message:err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category successfully updated",
		Data:    updatedCategory,
	})
}

func (cpc CategoryProductController) Delete(ctx *gin.Context){
	id := ctx.Param("id")
	categoryId,_ := strconv.Atoi(id)
	err := respository.DeleteCategory(cpc.Pool, categoryId)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success deleting category",
	})
}