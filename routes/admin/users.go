package admin

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UsersRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	userController := controller.UserController{Pool: pool}
	users := r.Group("/users")
	{
		users.GET("/", userController.GetUsers)
		users.DELETE("/:id", userController.DeleteUser)
		users.PATCH("/role/:id", userController.UpdateRole)
	}
	}
	