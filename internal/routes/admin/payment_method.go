package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func PaymentMethodRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	paymentMethodController := controller.PaymentMethodController{Pool: pool}

	paymentMethods := r.Group("/payment-methods")
	{
		paymentMethods.POST("", paymentMethodController.Create)
		paymentMethods.GET("", paymentMethodController.GetAll)
		paymentMethods.GET("/:id", paymentMethodController.GetByID)
		paymentMethods.PATCH("/:id", paymentMethodController.Edit)
		paymentMethods.DELETE("/:id", paymentMethodController.Delete)
	}
}
