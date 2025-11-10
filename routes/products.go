package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func ProductRoutes(r *gin.RouterGroup, conn *pgx.Conn){
	productController := controller.ProductController{Conn: conn}

	products := r.Group("/products")
	{
		products.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), productController.GetProducts)
		products.POST("/", middelware.VerifToken(), middelware.VerifRole("admin"), productController.CreateProduct)
		products.PATCH("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), productController.EditProduct)
		products.DELETE("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), productController.DeleteProduct)
		products.POST("/image/:id", productController.CreateImageProduct)
	}
}