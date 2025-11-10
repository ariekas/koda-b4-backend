package middelware

import (
	"back-end-coffeShop/lib/config"

	"github.com/gin-gonic/gin"
)

func CrossMiddelware(ctx *gin.Context){
	url := config.ReadEnvUrl()

	ctx.Header("AAccess-Control-Allow-Origin", url)
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
}