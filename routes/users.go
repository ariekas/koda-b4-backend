package routes

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UsersRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	userController := controller.UserController{Pool: pool}
	users := r.Group("/users")
	{
		users.GET("/", middelware.VerifToken(), middelware.VerifRole("admin"), userController.GetUsers)
		users.DELETE("/:id", middelware.VerifToken(), middelware.VerifRole("admin"), userController.DeleteUser)
		users.PATCH("/role/:id", userController.UpdateRole)
	}
	}
	