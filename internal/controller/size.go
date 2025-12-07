package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SizeController struct {
	Pool *pgxpool.Pool
}

func (sc SizeController) Create(ctx *gin.Context) {
	var req models.Size

	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("failed type request body")
	}

	size, err := respository.CreateSize(&req, sc.Pool)

	if err != nil {
		ctx.JSON(401, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(201, models.Response{
		Success: true,
		Message: "success created size",
		Result:  size,
	})
}

func (sc SizeController) GetAll(ctx *gin.Context) {
	sizes, err := respository.ListSize(sc.Pool)
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to get sizes",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Sizes retrieved successfully",
		"data":    sizes,
	})
}

func (sc SizeController) GetById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	size, err := respository.DetailSize(sc.Pool, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "Size not found",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(400, gin.H{
			"success": false,
			"message": "Failed to get size",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Size retrieved successfully",
		"data":    size,
	})
}

func (sc SizeController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	var req models.Size
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	size, err := respository.UpdateSize(sc.Pool, id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "Size not found",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to update size",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Size updated successfully",
		"data":    size,
	})
}

func (sc SizeController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(401, gin.H{
			"success": false,
			"message": "Invalid ID parameter",
			"error":   err.Error(),
		})
		return
	}

	err = respository.DeleteSize(sc.Pool, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(404, gin.H{
				"success": false,
				"message": "Size not found",
				"error":   err.Error(),
			})
			return
		}

		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Failed to delete size",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(201, gin.H{
		"success": true,
		"message": "Size deleted successfully",
	})
}
