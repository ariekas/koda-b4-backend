package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func OrderRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	OrdersController := controller.OrdersController{Pool: pool}

	orders := r.Group("/orders")
	{
		orders.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), OrdersController.GetOrders)
		orders.GET("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), OrdersController.GetById)
		orders.PATCH("/status/:id",middelware.VerifToken(), middelware.VerifRole("admin"), OrdersController.UpdateStatus)
	}
}