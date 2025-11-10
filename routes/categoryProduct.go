package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CategoryProductRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	categoryProductController := controller.CategoryProductController{Pool: pool}

	category := r.Group("/categorys")
	{
		category.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), categoryProductController.GetAll)
		category.GET("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), categoryProductController.GetByID)
		category.POST("/", middelware.VerifToken(), middelware.VerifRole("admin"), categoryProductController.Create)
		category.PATCH("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), categoryProductController.Edit)
		category.DELETE("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), categoryProductController.Delete)
	}
}