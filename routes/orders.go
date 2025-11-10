package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func OrderRoutes(r *gin.RouterGroup, conn *pgx.Conn){
	OrdersController := controller.OrdersController{Conn: conn}

	orders := r.Group("/orders")
	{
		orders.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), OrdersController.GetOrders)
		orders.PATCH("/status/:id",middelware.VerifToken(), middelware.VerifRole("admin"), OrdersController.UpdateStatus)
	}
}