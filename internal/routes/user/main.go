package user

import (
	"back-end-coffeShop/internal/controller"
	"back-end-coffeShop/internal/middelware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(r *gin.RouterGroup, pool *pgxpool.Pool) {
	productController := controller.ProductController{Pool: pool}
	cartController := controller.CartController{Pool: pool}
	transactionsControlelr := controller.TransactionsController{Pool: pool}
	usersController := controller.UserController{Pool: pool}
	historyController := controller.HistoryController{Pool: pool}
	categoryProductController := controller.CategoryProductController{Pool: pool}
	RatingController := controller.RatingController{Pool: pool}

	r.GET("/products", productController.GetProducts)
	r.GET("/products/favorite", productController.GetFavoriteProducts)
	r.GET("/products/filter", productController.Filter)
	r.GET("/product/:id", productController.DetailProduct)

	r.GET("/categories", categoryProductController.GetAll)

	r.POST("/cart", cartController.AddCart)
	r.GET("/cart", cartController.GetCart)
	r.GET("/cart/count", cartController.CountCart)
	r.DELETE("/cart/:id", cartController.DeleteCart)

	r.GET("/payment-method", transactionsControlelr.GetPaymentMethods)

	r.POST("/transactions", transactionsControlelr.CreateTransaction)

	r.GET("/user", usersController.GetUserLogin)
	r.PATCH("user/profile", usersController.UpdateProfile)

	r.GET("/historys", historyController.GetHistorys)
	r.GET("/history/:id", historyController.DetailHistory)

	r.POST("rating", middelware.VerifToken(), RatingController.CreateRating)
	r.GET("/ratings", middelware.VerifToken(), RatingController.GetMyRatings)
	r.GET("/rating/unrated/:transaction_id", middelware.VerifToken(), RatingController.GetUnratedProducts)
	r.PATCH("/rating/:id", middelware.VerifToken(), RatingController.UpdateRating)
	r.DELETE("/rating/:id", middelware.VerifToken(), RatingController.DeleteRating)
}
