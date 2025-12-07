package controller

import (
	"back-end-coffeShop/internal/config"
	"back-end-coffeShop/internal/models"
	"back-end-coffeShop/internal/respository"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentMethodController struct {
	Pool *pgxpool.Pool
}

var pathPaymentMethod = "/admin/payment-methods*"

func (pmc PaymentMethodController) GetAll(ctx *gin.Context) {
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
		response, err = respository.GetPaymentMethods(pmc.Pool, page)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		dataPaymentMethods, err := json.Marshal(response)
		if err == nil {
			config.Redis().Set(context.Background(), ctx.Request.RequestURI, dataPaymentMethods, 15*time.Minute)
		}
	} else {
		_ = json.Unmarshal([]byte(cache), &response)
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success getting data payment methods",
		Result:  response,
	})
}

func (pmc PaymentMethodController) Create(ctx *gin.Context) {
	var input models.PaymentMethodRequest
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Name is required",
		})
		return
	}

	file, err := ctx.FormFile("image_payment")
	var imageURL string

	if err == nil && file != nil {
		imageURL, err = config.UploaderFile(file, "payment_methods", fmt.Sprintf("payment_%d", time.Now().Unix()))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: fmt.Sprintf("Failed to upload image: %v", err),
			})
			return
		}
	}

	paymentMethod := models.PaymentMethods{
		Name:         input.Name,
		ImagePayment: imageURL,
	}

	result, err := respository.CreatePaymentMethod(pmc.Pool, paymentMethod)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(pathPaymentMethod)

	ctx.JSON(http.StatusCreated, models.Response{
		Success: true,
		Message: "Success create payment method",
		Result:  result,
	})
}

func (pmc PaymentMethodController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	paymentMethodId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	paymentMethod, err := respository.GetPaymentMethodById(pmc.Pool, paymentMethodId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success get data payment method",
		Result:  paymentMethod,
	})
}

func (pmc PaymentMethodController) Edit(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	var input models.PaymentMethodRequest
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	existingPM, err := respository.GetPaymentMethodById(pmc.Pool, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	file, err := ctx.FormFile("image_payment")
	imageURL := existingPM.ImagePayment 

	if err == nil && file != nil {
		imageURL, err = config.UploaderFile(file, "payment_methods", fmt.Sprintf("payment_%d", time.Now().Unix()))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: fmt.Sprintf("Failed to upload image: %v", err),
			})
			return
		}
	}

	paymentMethod := models.PaymentMethods{
		Name:         input.Name,
		ImagePayment: imageURL,
	}

	updatedPaymentMethod, err := respository.EditPaymentMethod(pmc.Pool, id, paymentMethod)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(pathPaymentMethod)

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Payment method successfully updated",
		Result:  updatedPaymentMethod,
	})
}

func (pmc PaymentMethodController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	paymentMethodId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "Invalid ID",
		})
		return
	}

	err = respository.DeletePaymentMethod(pmc.Pool, paymentMethodId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	config.InvalidateRedis(pathPaymentMethod)

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success deleting payment method",
	})
}
