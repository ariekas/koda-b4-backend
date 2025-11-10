package main

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/routes"
	"fmt"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "back-end-coffeShop/docs" 
	"github.com/gin-gonic/gin"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5..."
func main() {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	connectDb := controller.ConnectDB()
	r := gin.Default()

	routes.UsersRoutes(&r.RouterGroup, connectDb)
	routes.AuthRoutes(&r.RouterGroup, connectDb)
	routes.ProductRoutes(&r.RouterGroup, connectDb)
	routes.OrderRoutes(&r.RouterGroup, connectDb)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")

}