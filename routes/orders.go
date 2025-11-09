package routes

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func OrderRoutes(r *gin.RouterGroup, conn *pgx.Conn){
	OrdersController := controller.OrdersController{Conn: conn}

	orders := r.Group("/orders")
	{
		orders.GET("/", OrdersController.GetOrders)
		orders.PATCH("/update/status/:id", OrdersController.UpdateStatus)
	}
}