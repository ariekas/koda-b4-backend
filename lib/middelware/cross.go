package middelware

import (
	"back-end-coffeShop/lib/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CrossMiddelware() gin.HandlerFunc {
	url := config.ReadEnvUrl()
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", url}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE","PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	return cors.New(config)
}