package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func VariantRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	variantController := controller.VariantController{Pool: pool}
	variants := r.Group("/variants")
	{
		variants.POST("", variantController.Create)
		variants.GET("", variantController.List)
		variants.GET("/:id", variantController.Detail)
		variants.PATCH("/:id", variantController.Update)
		variants.DELETE("/:id", variantController.Delete)
	}
}
