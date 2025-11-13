package handler

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/models"
	"back-end-coffeShop/routes"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	App  *gin.Engine
	pool *pgxpool.Pool
)

func init() {
	App = gin.New()
	App.RedirectTrailingSlash = false
	App.RedirectFixedPath = false
	App.Use(gin.Recovery(), gin.Logger())
	pool = controller.ConnectDB()

	App = gin.New()
	App.Use(gin.Recovery(), gin.Logger())

	App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "Backend is running well 🚀",
		})
	})

	routes.MainRoutes(App, pool)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if err := pool.Ping(context.Background()); err != nil {
		fmt.Println("Reconnecting to database...")
		controller.ConnectDB()

	}

	App.ServeHTTP(w, r)
}
