package controller

import (
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeliveryController struct {
	Pool *pgxpool.Pool
}

func (dc *DeliveryController) CreateDelivery(ctx *gin.Context) {
	var req models.CreateDeliveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery := &models.Deliverys{
		Name:  req.Name,
		Price: req.Price,
	}

	if err := respository.CreateDelivery(dc.Pool, delivery); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Delivery created successfully",
		"data":    delivery,
	})
}

func (dc *DeliveryController) GetAllDeliveries(ctx *gin.Context) {
	deliveries, err := respository.ListDelivery(dc.Pool)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Deliveries retrieved successfully",
		"data":    deliveries,
	})
}

func (dc *DeliveryController) GetDeliveryByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	delivery, err := respository.DetailDelivery(dc.Pool, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delivery retrieved successfully",
		"data":    delivery,
	})
}


func (dc *DeliveryController) UpdateDelivery(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	var req models.UpdateDeliveryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delivery := &models.Deliverys{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
	}

	if err := respository.UpdateDelivery(dc.Pool, delivery); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delivery updated successfully",
		"data":    delivery,
	})
}

func (dc *DeliveryController) DeleteDelivery(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery ID"})
		return
	}

	if err := respository.DeleteDelivery(dc.Pool, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Delivery deleted successfully",
	})
}