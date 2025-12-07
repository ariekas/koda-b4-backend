package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TaxRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	taxController := controller.TaxController{Pool: pool}

	taxes := r.Group("/taxes")
	{
		taxes.POST("", taxController.CreateTax)
		taxes.GET("", taxController.GetAllTaxes)
		taxes.GET("/:id", taxController.GetTaxByID)
		taxes.PATCH("/:id", taxController.UpdateTax)
		taxes.DELETE("/:id", taxController.DeleteTax)
	}
}
