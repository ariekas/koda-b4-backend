package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SizeRoutes(router *gin.RouterGroup, pool *pgxpool.Pool) {
	sizeController := controller.SizeController{Pool: pool}
	sizes := router.Group("/sizes")
	{
		sizes.POST("", sizeController.Create)
		sizes.GET("", sizeController.GetAll)
		sizes.GET("/:id", sizeController.GetById)
		sizes.PATCH("/:id", sizeController.Update)
		sizes.DELETE("/:id", sizeController.Delete)
	}
}
