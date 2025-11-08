package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func UsersRoutes(r *gin.RouterGroup, conn *pgx.Conn) {
	userController := controller.UserController{Conn: conn}
	users := r.Group("/users")
	{
		users.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), userController.GetUsers)
	}
	}
	