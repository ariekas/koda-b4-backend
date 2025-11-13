package user

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	productController := controller.ProductController{Pool: pool}
	cartController := controller.CartController{Pool: pool}
	transactionsControlelr :=  controller.TransactionsController{Pool: pool}
	usersController := controller.UserController{Pool: pool}
	historyController := controller.HistoryController{Pool: pool}

	r.GET("/products/favorite", productController.GetFavoriteProducts)
	r.GET("/products/filter", productController.Filter)
	r.GET("/product/:id", productController.DetailProduct)
	r.POST("/cart", cartController.AddCart)
	r.GET("/cart", cartController.GetCart)
	r.POST("/transactions",transactionsControlelr.CreateTransaction)
	r.GET("/users/:id", usersController.GetUserById)
	r.PATCH("users/profile/:id", usersController.UpdateProfile)
	r.GET("/historys", historyController.GetHistorys)
	r.GET("/history/:id", historyController.DetailHistory)
}