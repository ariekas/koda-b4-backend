package main

import (
	"back-end-coffeShop/internal/controller"
	"back-end-coffeShop/internal/middelware"
	"back-end-coffeShop/internal/routes"
	"fmt"

	_ "back-end-coffeShop/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	connectDb := controller.ConnectDB()
	r := gin.Default()

	r.MaxMultipartMemory = 8 << 20
	r.Use(middelware.CrossMiddelware())
	r.Use(middelware.AllowPreflight)
	routes.MainRoutes(r, connectDb)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")

}