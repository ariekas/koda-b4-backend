package routes

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func AuthRoutes(r *gin.RouterGroup, conn *pgx.Conn){
	authController := controller.AuthRegister{Conn: conn}
	auth := r.Group("/")
	{
		auth.POST("/register", authController.Register)
	}
}