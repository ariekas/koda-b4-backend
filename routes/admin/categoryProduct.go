package admin

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryProductRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	categoryProductController := controller.CategoryProductController{Pool: pool}

	category := r.Group("/categorys")
	{
		category.GET("/",  categoryProductController.GetAll)
		category.GET("/:id",  categoryProductController.GetByID)
		category.POST("/",  categoryProductController.Create)
		category.PATCH("/:id",  categoryProductController.Edit)
		category.DELETE("/:id",  categoryProductController.Delete)
	}
}