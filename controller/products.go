package controller

import (
	"back-end-coffeShop/lib/config"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
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

// GetProducts godoc
// @Summary Get all products
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Response "Success getting products"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /products [get]
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

	var response respository.PaginationResponse

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
		Data:    response,
	})
}


// CreateProduct godoc
// @Summary Create a new product
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body models.ProductInput true "Product data"
// @Success 201 {object} models.Response "Success create product"
// @Failure 400 {object} models.Response "Invalid input"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /products [post]
func (pc ProductController) CreateProduct(ctx *gin.Context) {
	product := respository.Create(ctx, pc.Pool)

	redis := config.Redis()
	iter := redis.Scan(context.Background(), 0 ,"/products*", 0).Iterator()
	for iter.Next(context.Background()) {
		redis.Del(context.Background(), iter.Val())
	}
	if err := iter.Err(); err != nil {
		fmt.Println("Redis scan error:", err)
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success Create product",
		Data:    product,
	})
}

// EditProduct godoc
// @Summary Edit product
// @Tags Products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.ProductInput true "Updated product data"
// @Success 200 {object} models.Response "Success edit product"
// @Failure 404 {object} models.Response "Product not found"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /products/{id} [patch]
func (pc ProductController) EditProduct(ctx *gin.Context) {
	newProduct, err := respository.Edit(pc.Pool, ctx)

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
		Data:    newProduct,
	})
}

// DeleteProduct godoc
// @Summary Delete product
// @Tags Products
// @Security BearerAuth
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Response "Success delete product"
// @Failure 404 {object} models.Response "Product not found"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /products/{id} [delete]
func (pc ProductController) DeleteProduct(ctx *gin.Context) {
	err := respository.Delete(pc.Pool, ctx)

	if err != nil {
		ctx.JSON(404, models.Response{
			Success: false,
			Message: "Error: Failed to get product",
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Sucess deleted product",
	})
}

// CreateImageProduct godoc
// @Summary Upload product images
// @Tags Products
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Product ID"
// @Param images formData file true "Upload product images"
// @Success 201 {object} models.Response "Success create image product"
// @Failure 400 {object} models.Response "Failed to upload image"
// @Failure 401 {object} models.Response "Unauthorized"
// @Router /products/image/{id} [post]
func (pc ProductController) CreateImageProduct(ctx *gin.Context) {
	id := ctx.Param("id")
	productId, _ := strconv.Atoi(id)

	from, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "Failed to read form data",
		})
		return
	}

	files := from.File["images"]
	if len(files) == 0 {
		ctx.JSON(400, models.Response{
			Success: false,
			Message: "No image uploaded",
		})
		return
	}

	imageProduct, err := respository.CreateImageProduct(pc.Pool, ctx, productId, files)
	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: "Error: Failed to create product image",
		})
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success create image product",
		Data:    imageProduct,
	})
}

func (pc ProductController) GetAllImageProduct(ctx *gin.Context) {
	images, err := respository.GetAllImageProduct(pc.Pool)

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
		Data:    images,
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

	response := respository.PaginationResponse{
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
		Data:  response,
	})
}

func (pc ProductController) Filter(ctx *gin.Context) {
	name := ctx.Query("name")
	categoryStr := ctx.Query("category")
	sortBy := ctx.Query("sort_by")
	priceMin := ctx.Query("price_min")
	priceMax := ctx.Query("price_max")
	pageQuery := ctx.Query("page")
	limitQuery := ctx.Query("limit")

	page := 1
	limit := 20
	if p, err := strconv.Atoi(pageQuery); err == nil && p > 0 {
		page = p
	}
	if l, err := strconv.Atoi(limitQuery); err == nil && l > 0 {
		limit = l
	}

	data, total, err := respository.FilterProducts(pc.Pool, name, categoryStr, sortBy, priceMin, priceMax, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Error: %v", err),
		})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	baseURL := fmt.Sprintf("/products/filter?name=%s&category=%s&sort_by=%s&price_min=%s&price_max=%s&limit=%d",
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

	response := respository.PaginationResponse{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		Links:      links,
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success filtering products",
		Data:    response,
	})
}