package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VariantController struct {
	Pool *pgxpool.Pool
}

func (vc VariantController) Create(ctx *gin.Context) {
	var req models.Variant

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	variant, err := respository.CreateVariant(vc.Pool, &req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":   "Failed to create variant product",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Variant product created successfully",
		"data":    variant,
	})
}

func (vc VariantController) List(ctx *gin.Context) {
	variants, err := respository.ListVariant(vc.Pool)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error":   "Failed to get variant products",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Variant products retrieved successfully",
		"data":    variants,
	})
}

func (vc VariantController) Detail(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	variant, err := respository.DetailVariant(vc.Pool, id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"error":   "Variant product not found",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Variant product retrieved successfully",
		"data":    variant,
	})
}

func (vc VariantController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	var req models.Variant
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
		return
	}

	variant, err := respository.UpdateVariant(vc.Pool, id, &req)
	if err != nil {
		if err.Error() == "variant product not found" {
			ctx.JSON(404, gin.H{
				"error":   "Variant product not found",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(500, gin.H{
			"error":   "Failed to update variant product",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Variant product updated successfully",
		"data":    variant,
	})
}

func (vc VariantController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a number",
		})
		return
	}

	err = respository.DeleteVariant(vc.Pool, id)
	if err != nil {
		if err.Error() == "variant product not found" {
			ctx.JSON(404, gin.H{
				"error":   "Variant product not found",
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(500, gin.H{
			"error":   "Failed to delete variant product",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"message": "Variant product deleted successfully",
	})
}
