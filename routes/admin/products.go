package admin


import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	productController := controller.ProductController{Pool: pool}

	products := r.Group("/products")
	{
		products.GET("/",  productController.GetProducts)
		products.POST("/",  productController.CreateProduct)
		products.PATCH("/:id",  productController.EditProduct)
		products.DELETE("/:id",  productController.DeleteProduct)
		products.POST("/image/:id", productController.CreateImageProduct)
		products.GET("/images", productController.GetAllImageProduct)
		products.DELETE("/image/:id", productController.DeleteImageProduct)
	}
}