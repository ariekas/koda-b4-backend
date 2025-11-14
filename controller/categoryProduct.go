package controller

import (
	"back-end-coffeShop/lib/config"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryProductController struct {
	Pool *pgxpool.Pool
}

var path  = "/categorys*"

func (cpc CategoryProductController) GetAll(ctx *gin.Context) {
	cache, err := config.Redis().Get(context.Background(), ctx.Request.RequestURI).Result()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "error: Failed to connect redis",
		})
	}

	pageQuery := ctx.Query("page")
	page := 1

	if pageQuery != "" {
		if p, err := strconv.Atoi(pageQuery); err == nil && p > 0 {
			page = p
		}
	}

	var response models.PaginationResponse

	if cache == "" {
		response, err := respository.GetCategories(cpc.Pool, page)

		if err != nil {
			ctx.JSON(401, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		dataCategorys, err := json.Marshal(response)
		if err == nil {
			config.Redis().Set(context.Background(), ctx.Request.RequestURI, dataCategorys, 15*time.Minute)
		}
	} else {
		_ = json.Unmarshal([]byte(cache), &response)
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success getting data category",
		Data:    response,
	})
}

func (cpc CategoryProductController) Create(ctx *gin.Context) {

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

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(path)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success create category",
		Data:    categorys,
	})
}

func (cpc CategoryProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryId, _ := strconv.Atoi(id)
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
		Data:    category,
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
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(path)

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Category successfully updated",
		Data:    updatedCategory,
	})
}

func (cpc CategoryProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	categoryId, _ := strconv.Atoi(id)
	err := respository.DeleteCategory(cpc.Pool, categoryId)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(path)
	
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success deleting category",
	})
}
