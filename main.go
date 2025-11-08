package main

import (
	"back-end-coffeShop/controller"
	"back-end-coffeShop/routes"
	"fmt"

	"github.com/gin-gonic/gin"
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

	routes.UsersRoutes(&r.RouterGroup, connectDb)
	routes.AuthRoutes(&r.RouterGroup, connectDb)
	routes.ProductRoutes(&r.RouterGroup, connectDb)

	r.Run(":8080")

}