package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductController struct{
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
func (pc ProductController) GetProducts(ctx *gin.Context){
	produts, err := respository.GetProducts(pc.Pool)

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
func (pc ProductController) CreateProduct(ctx *gin.Context){
	product := respository.Create(ctx, pc.Pool)

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "Success Create product",
		Data: product,
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
func (pc ProductController) EditProduct(ctx *gin.Context){
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
		Data: newProduct,
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
func (pc ProductController) DeleteProduct(ctx *gin.Context){
	err := respository.Delete(pc.Pool, ctx)

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
func (pc ProductController) CreateImageProduct(ctx *gin.Context){
	id := ctx.Param("id")
	productId, _ := strconv.Atoi(id)

	from , err := ctx.MultipartForm()
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
		Data: imageProduct,
	})
}