package handler

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/models"
	"back-end-coffeShop/routes"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var App *gin.Engine

func InitHandler(pool *pgxpool.Pool) {
	App = gin.New()
	App.Use(gin.Recovery())

	App.Use(gin.Logger())

	App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Backend is running well 🚀",
		})
	})

	routes.MainRoutes(App, pool)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	pool := controller.ConnectDB()

	if App == nil {
		InitHandler(pool)
	}

	App.ServeHTTP(w, r)
}
