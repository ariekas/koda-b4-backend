package middelware

import (
	"back-end-coffeShop/lib/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifToken() gin.HandlerFunc {
	jwtToken := config.ReadENV()

	return func(ctx *gin.Context)  {
		authHeader := ctx.Request.Header.Get("Authorization")

		tokenString,_ := strings.CutPrefix(authHeader, "Bearer ")
		
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