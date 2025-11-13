package api

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/routes"
	"fmt"
	"net/http"

	_ "back-end-coffeShop/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var App *gin.Engine

func init() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	connectDb := controller.ConnectDB()
	App = gin.Default()

	App.MaxMultipartMemory = 8 << 20
	App.Use(middelware.CrossMiddelware)
	App.Use(middelware.AllowPreflight)
	
	routes.MainRoutes(App, connectDb)

	App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func Handler(w http.ResponseWriter, r *http.Request) {
	App.ServeHTTP(w, r)
}