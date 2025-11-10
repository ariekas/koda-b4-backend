package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	productController := controller.ProductController{Pool: pool}

	products := r.Group("/products")
	{
		products.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), productController.GetProducts)
		products.POST("/", middelware.VerifToken(), middelware.VerifRole("admin"), productController.CreateProduct)
		products.PATCH("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), productController.EditProduct)
		products.DELETE("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), productController.DeleteProduct)
		products.POST("/image/:id", productController.CreateImageProduct)
		products.GET("/images", productController.GetAllImageProduct)
		products.DELETE("/image/:id", productController.DeleteImageProduct)
	}
}