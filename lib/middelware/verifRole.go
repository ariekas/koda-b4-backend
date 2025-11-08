package middelware

import (
	"back-end-coffeShop/models"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func VerifRole(roles ...string) gin.HandlerFunc {
	godotenv.Load()

	var jwtToken = os.Getenv("JWT_TOKEN")

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

		claims,_ := token.Claims.(jwt.MapClaims)

		userRole,_ := claims["role"].(string)

		allowed := false

		for _, role := range roles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(400, models.Response{
				Success: false,
				Message: "Error: Access role Denied",
			})
			ctx.Abort()
			return 
		}
		ctx.Next()	
	}
}