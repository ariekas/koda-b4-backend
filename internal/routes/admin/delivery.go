package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DeliveryRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	deliveryController := controller.DeliveryController{Pool: pool}

	deliveries := r.Group("/deliveries")
	{
		deliveries.POST("", deliveryController.CreateDelivery)
		deliveries.GET("", deliveryController.GetAllDeliveries)
		deliveries.GET("/:id", deliveryController.GetDeliveryByID)
		deliveries.PATCH("/:id", deliveryController.UpdateDelivery)
		deliveries.DELETE("/:id", deliveryController.DeleteDelivery)
	}
}
