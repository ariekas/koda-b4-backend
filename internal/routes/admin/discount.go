package admin

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func DiscountRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	discountController := controller.DiscountController{Pool: pool}

	discount := r.Group("/discounts")
	{
		discount.GET("", discountController.GetAll)       
		discount.POST("", discountController.Create)       
		discount.GET("/:id", discountController.Detail)    
		discount.PATCH("/:id", discountController.Update)  
		discount.DELETE("/:id", discountController.Delete) 
	}
}
