package admin

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrderRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	OrdersController := controller.OrdersController{Pool: pool}

	orders := r.Group("/orders")
	{
		orders.GET("/",  OrdersController.GetOrders)
		orders.GET("/:id",  OrdersController.GetById)
		orders.PATCH("/status/:id", OrdersController.UpdateStatus)
	}
}