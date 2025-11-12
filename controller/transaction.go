package controller

import (
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsController struct {
	Pool *pgxpool.Pool
}

// GetOrders godoc
// @Summary Get all orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /orders [get]
func (tc TransactionsController) GetTransactions(ctx *gin.Context) {
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

	response, err := respository.GetTransactions(tc.Pool, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Error getting transactions: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success getting transactions",
		Data:    response,
	})
}

func (tc TransactionsController) GetTransactionById(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid transaction ID",
		})
		return
	}

	transaction, err := respository.GetTransactionById(tc.Pool, transactionId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: "Transaction not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success getting transaction detail",
		Data:    transaction,
	})
}

// UpdateStatus godoc
// @Summary Update status order
// @Tags Orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body models.InputNewStatus true "New Status Data"
// @Success 200 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /orders/status/{id} [patch]
func (tc TransactionsController) UpdateTransactionStatus(ctx *gin.Context) {
	id := ctx.Param("id")
	transactionId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid transaction ID",
		})
		return
	}

	var input models.InputNewStatus
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: fmt.Sprintf("Invalid input: %v", err),
		})
		return
	}

	if err := respository.UpdateTransactionStatus(tc.Pool, transactionId, input.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to update transaction status: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success updating transaction status",
	})
}
