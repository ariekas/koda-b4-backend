package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaxController struct {
	Pool *pgxpool.Pool
}

func (tc *TaxController) CreateTax(ctx *gin.Context) {
	var req models.CreateTaxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tax := &models.CreateTaxRequest{
		Name: req.Name,
		Tax:  req.Tax,
	}

	if err := respository.CreateTax(tc.Pool, tax); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Tax created successfully",
		"data":    tax,
	})
}

func (tc *TaxController) GetAllTaxes(ctx *gin.Context) {
	taxes, err := respository.ListTax(tc.Pool)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Taxes retrieved successfully",
		"data":    taxes,
	})
}

func (tc *TaxController) GetTaxByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tax ID"})
		return
	}

	tax, err := respository.DetailTax(tc.Pool, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tax retrieved successfully",
		"data":    tax,
	})
}

func (tc *TaxController) UpdateTax(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tax ID"})
		return
	}

	var req models.UpdateTaxRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tax := &models.Taxe{
		ID:   id,
		Name: req.Name,
		Tax:  req.Tax,
	}

	if err := respository.UpdateTax(tc.Pool, tax); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tax updated successfully",
		"data":    tax,
	})
}
func (tc *TaxController) DeleteTax(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tax ID"})
		return
	}

	if err := respository.DeleteTax(tc.Pool, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tax deleted successfully",
	})
}
