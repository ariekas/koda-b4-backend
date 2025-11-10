package main

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/routes"
	"fmt"

	_ "back-end-coffeShop/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.MaxMultipartMemory = 8 << 20
	r.Use(middelware.CrossMiddelware)
	r.Use(middelware.AllowPreflight)
	
	routes.UsersRoutes(&r.RouterGroup, connectDb)
	routes.AuthRoutes(&r.RouterGroup, connectDb)
	routes.ProductRoutes(&r.RouterGroup, connectDb)
	routes.OrderRoutes(&r.RouterGroup, connectDb)
	routes.CategoryProductRoutes(&r.RouterGroup, connectDb)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")

}