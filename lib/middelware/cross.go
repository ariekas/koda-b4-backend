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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	}

	return cors.New(config)
}