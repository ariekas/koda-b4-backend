package routes

import (
	"back-end-coffeShop/internal/middelware"
	routesAdmin "back-end-coffeShop/internal/routes/admin"
	routesAuth "back-end-coffeShop/internal/routes/auth"
	routesUser "back-end-coffeShop/internal/routes/user"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MainRoutes(r *gin.Engine, pool *pgxpool.Pool) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "success",
			"message": "Backend berjalan",
		})
	})

	admin := r.Group("/admin")
	admin.Use(middelware.VerifToken(), middelware.VerifRole("admin"))
	{
		routesAdmin.UsersRoutes(admin, pool)
		routesAdmin.CategoryProductRoutes(admin, pool)
		routesAdmin.ProductRoutes(admin, pool)
		routesAdmin.TransactionRoutes(admin, pool)
		routesAdmin.DiscountRoutes(admin, pool)
		routesAdmin.SizeRoutes(admin, pool)
		routesAdmin.VariantRoutes(admin, pool)
		routesAdmin.TaxRoutes(admin, pool)
		routesAdmin.DeliveryRoutes(admin, pool)
		routesAdmin.PaymentMethodRoutes(admin, pool)
	}

	auth := r.Group("/auth")
	{
		routesAuth.AuthRoutes(auth, pool)
	}

	user := r.Group("/")
	{
		routesUser.UserRoutes(user, pool)
	}
}

