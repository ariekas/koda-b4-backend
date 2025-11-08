package routes

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func ProductRoutes(r *gin.RouterGroup, conn *pgx.Conn){
	productController := controller.ProductController{Conn: conn}

	products := r.Group("/products")
	{
		products.GET("/", productController.GetProducts)
		products.POST("/", productController.CreateProduct)
		products.PATCH("/edit/:id", productController.EditProduct)
	}
}