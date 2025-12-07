package controller

import (
	"back-end-coffeShop/internal/config"
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductController struct {
	Pool *pgxpool.Pool
}

var Path = "/admin/products*"

func (pc ProductController) GetProducts(ctx *gin.Context) {
	cache, err := config.Redis().Get(context.Background(), ctx.Request.RequestURI).Result()
	if err != nil {
		fmt.Println("Error: Redis", err)
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
		response, err = respository.GetProducts(pc.Pool, page)
		if err != nil {
			ctx.JSON(500, models.Response{
				Success: false,
				Message: "Failed getting data products",
			})
			return
		}

		dataProduct, err := json.Marshal(response)
		if err == nil {
			config.Redis().Set(context.Background(), ctx.Request.RequestURI, dataProduct, 5*time.Minute)
		}

	} else {
		_ = json.Unmarshal([]byte(cache), &response)
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success getting data products",
		Result:  response,
	})
}

func (pc ProductController) CreateProduct(ctx *gin.Context) {
	var input models.ProductCreateInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid input",
		})
		return
	}

	createdProduct, err := respository.CreateProduct(pc.Pool, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(Path)
	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success Create product",
		Result:  createdProduct,
	})
}

func (pc ProductController) EditProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid product ID",
		})
		return
	}

	var input models.ProductUpdateInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	if input.Name != nil && *input.Name == "" {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Product name cannot be empty",
		})
		return
	}

	if input.Price != nil && *input.Price < 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Price cannot be negative",
		})
		return
	}

	if input.Stock != nil && *input.Stock < 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Stock cannot be negative",
		})
		return
	}

	err = respository.EditProduct(pc.Pool, productId, input)

	if err != nil {
		if err.Error() == "product not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "Product not found",
				Result:  nil,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to edit product: " + err.Error(),
			Result:  nil,
		})
		return
	}

	config.InvalidateRedis(Path)

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Product updated successfully",
	})
}

func (pc ProductController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	productId, _ := strconv.Atoi(id)
	err := respository.DeleteProduct(pc.Pool, productId)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(Path)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Sucess deleted product",
	})
}

func (pc ProductController) CreateImageProduct(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("id"))

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: "Cannot read form data"})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		ctx.JSON(400, models.Response{Success: false, Message: "No image uploaded"})
		return
	}

	images, err := respository.CreateImageProduct(pc.Pool, ctx, productId, files)
	if err != nil {
		ctx.JSON(400, models.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success create image product",
		Result:  images,
	})
}

func (pc ProductController) GetAllImageProduct(ctx *gin.Context) {
	images, err := respository.GetAllImageProducts(pc.Pool)

	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Failed to get image products",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get all image products",
		Result:  images,
	})
}

func (pc ProductController) DeleteImageProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	imageId, _ := strconv.Atoi(id)

	err := respository.DeleteImageProduct(pc.Pool, imageId)
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Failed delete image product",
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success delete image product",
	})
}

func (pc ProductController) GetFavoriteProducts(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	limit := 4

	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
		limit = l
	}

	data, total, err := respository.GetProductFavorite(pc.Pool, limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	response := models.PaginationResponse{
		Data:       data,
		Page:       1,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
		Links: map[string]string{
			"self": fmt.Sprintf("/products/favorite?limit=%d", limit),
		},
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success getting favorite product",
		Result:  response,
	})
}

func (pc ProductController) Filter(ctx *gin.Context) {
	name := ctx.Query("name")
	categoryStr := ctx.Query("category")

	sortBy := ctx.Query("sortBy")
	priceMin, _ := strconv.ParseFloat(ctx.Query("priceMin"), 64)
	priceMax, _ := strconv.ParseFloat(ctx.Query("priceMax"), 64)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))

	data, total, err := respository.FilterProducts(
		pc.Pool,
		name,
		categoryStr,
		sortBy,
		priceMin,
		priceMax,
		page,
		limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed filtering products: %v", err),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	if totalPages < 1 {
		totalPages = 1
	}

	baseURL := fmt.Sprintf(
		"/products/filter?name=%s&category=%s&sortBy=%s&priceMin=%v&priceMax=%v&limit=%d",
		name, categoryStr, sortBy, priceMin, priceMax, limit,
	)

	links := map[string]string{
		"self":  fmt.Sprintf("%s&page=%d", baseURL, page),
		"first": fmt.Sprintf("%s&page=1", baseURL),
		"last":  fmt.Sprintf("%s&page=%d", baseURL, totalPages),
	}

	if page > 1 {
		links["prev"] = fmt.Sprintf("%s&page=%d", baseURL, page-1)
	}
	if page < totalPages {
		links["next"] = fmt.Sprintf("%s&page=%d", baseURL, page+1)
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success filtering products",
		Result: models.PaginationResponse{
			Data:       data,
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			Links:      links,
		},
	})
}

func (pc ProductController) DetailProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	productId, _ := strconv.Atoi(id)
	product, err := respository.DetailProduct(pc.Pool, productId)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(200, models.Response{
		Success: true,
		Message: "Success get product detail",
		Result:  product,
	})
}
