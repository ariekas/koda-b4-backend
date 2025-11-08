package middelware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func VerifToken() gin.HandlerFunc {
	godotenv.Load()

	var jwtToken = os.Getenv("JWT_TOKEN")

	return func(ctx *gin.Context)  {
		authHeader := ctx.Request.Header.Get("Authorization")

		tokenString,_ := strings.CutPrefix(authHeader, "Bearer")
		
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(jwtToken), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"message": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()	

	}
}