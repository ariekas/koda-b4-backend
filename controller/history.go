package controller

import (
	"back-end-coffeShop/lib/middelware"
	"back-end-coffeShop/models"
	"back-end-coffeShop/respository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HistoryController struct{
	Pool *pgxpool.Pool
}

func (hc HistoryController) GetHistorys(ctx *gin.Context){
	userID := middelware.GetUserFromToken(ctx)
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Token tidak valid atau user tidak ditemukan",
		})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	month := ctx.Query("month") 
	status := ctx.Query("status")

	historys, err := respository.GetHistorys(hc.Pool, userID, page, limit, month, status)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "Failed to getting data history transaction",
		})
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success getting data history",
		Data: historys,
	})
}

func (hc HistoryController) DetailHistory(ctx *gin.Context){
	userId := middelware.GetUserFromToken(ctx)
	id := ctx.Param("id")
	historyId,_ := strconv.Atoi(id)
	if userId == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Token tidak valid atau user tidak ditemukan",
		})
		return
	}

	history, err := respository.DetailHistory(hc.Pool, userId,historyId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: "Failed get history not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "Success getting detail history",
		Data: history,
	})
}