package middelware

import (
	"back-end-coffeShop/internal/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)
func VerifToken() gin.HandlerFunc {
    jwtSecret := config.ReadENV()

    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            ctx.JSON(401, gin.H{"message": "Authorization header required"})
            ctx.Abort()
            return
        }

        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
            return []byte(jwtSecret), nil
        })

        if err != nil || !token.Valid {
            ctx.JSON(401, gin.H{"message": "Invalid token"})
            ctx.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            ctx.JSON(401, gin.H{"message": "Invalid token claims"})
            ctx.Abort()
            return
        }

        userIDFloat, ok := claims["user_id"].(float64)
        if !ok {
            ctx.JSON(401, gin.H{"message": "User ID not found in token"})
            ctx.Abort()
            return
        }

        userID := int(userIDFloat)

        ctx.Set("user_id", userID)

        ctx.Next()
    }
}
