package controller

import (
	"back-end-coffeShop/internal/middelware"
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionsController struct {
	Pool *pgxpool.Pool
}

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
		Result:  response,
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
		Result:  transaction,
	})
}

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

	if err := respository.UpdateTransactionStatus(tc.Pool, transactionId, input.StatusTransactionID); err != nil {
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

func (tc TransactionsController) CreateTransaction(ctx *gin.Context) {
	userID := middelware.GetUserFromToken(ctx)

	carts, err := respository.GetCartTransaction(tc.Pool, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to get cart: %v", err),
		})
		return
	}
	if len(carts) == 0 {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Cart is empty",
		})
		return
	}

	var input models.TransactionInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: fmt.Sprintf("Invalid input: %v", err),
		})
		return
	}

	subtotal := 0.0
	for _, item := range carts {
		priceAfterDiscount := item.ProductPrice - item.DiscountAmount
		itemTotal := (priceAfterDiscount + item.SizeCost + item.VariantCost) * float64(item.Quantity)
		subtotal += itemTotal
	}

	deliveryPrice, err := respository.GetDelivery(tc.Pool, input.DeliveryID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: fmt.Sprintf("Invalid delivery method: %v", err),
		})
		return
	}

	taxID := input.TaxID
	if taxID == 0 {
		taxID = 1
	}
	tax, err := respository.GetTax(tc.Pool, taxID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: fmt.Sprintf("Invalid tax: %v", err),
		})
		return
	}

	taxAmount := subtotal * (tax.Percentage / 100)

	total := subtotal + taxAmount + deliveryPrice

	invoice := fmt.Sprintf(
		"INV-%d-%s",
		userID,
		time.Now().Format("20060102"),
	)

	tx, err := tc.Pool.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to start transaction",
		})
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			tx.Commit(context.Background())
		}
	}()

	txID, err := respository.CreateTransaction(tc.Pool, userID, input, subtotal, taxAmount, total, invoice, tx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to create transaction: %v", err),
		})
		return
	}

	for _, item := range carts {
		if err := respository.CreateTransactionItem(tc.Pool, tx, txID, item); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: fmt.Sprintf("Failed to add transaction item: %v", err),
			})
			return
		}
	}

	if err := respository.ClearCart(tc.Pool, tx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to clear cart",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Transaction created successfully",
		Result: models.TransactionResponse{
			Invoice:       invoice,
			Total:         total,
			PaymentStatus: "pending",
		},
	})
}

func (tc TransactionsController) GetPaymentMethods(c *gin.Context) {
	data, err := respository.GetPaymentMethod(tc.Pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to fetch payment methods",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success get payment methods",
		Result:  data,
	})
}

func (tc TransactionsController) GetStatus(ctx *gin.Context) {
	statuses, err := respository.GetAllStatus(tc.Pool)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success get all status",
		Result:  statuses,
	})
}

func (tc TransactionsController) GetDeliveries(ctx *gin.Context) {
	deliveries, err := respository.GetAllDeliveries(tc.Pool)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to get deliveries: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success get all deliveries",
		Result:  deliveries,
	})
}

func (tc TransactionsController) GetTaxes(ctx *gin.Context) {
	taxes, err := respository.GetAllTaxes(tc.Pool)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: fmt.Sprintf("Failed to get taxes: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success get all taxes",
		Result:  taxes,
	})
}
