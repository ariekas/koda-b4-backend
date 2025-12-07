package auth

import (
	"back-end-coffeShop/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func AuthRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	authController := controller.AuthController{Pool: pool}
	auth := r.Group("/")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)
		auth.POST("/forget-password", authController.ForgetPassword)
		auth.PATCH("/reset-password", authController.ResetPassword)
	}
}
