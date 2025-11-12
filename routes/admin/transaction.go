package admin

import (
	"back-end-coffeShop/controller"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TransactionRoutes(r *gin.RouterGroup, pool *pgxpool.Pool){
	TransactionsController := controller.TransactionsController{Pool: pool}

	transaction := r.Group("/transaction")
	{
		transaction.GET("/",  TransactionsController.GetTransactions)
		transaction.GET("/:id",  TransactionsController.GetTransactionById)
		transaction.PATCH("/status/:id", TransactionsController.UpdateTransactionStatus)
	}
}