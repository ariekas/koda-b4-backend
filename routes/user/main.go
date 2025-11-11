package user

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	productController := controller.ProductController{Pool: pool}

	r.GET("/products/favorite", productController.GetFavoriteProducts)
	r.GET("/products/filter", productController.Filter)
}