package middelware

import (
	"back-end-coffeShop/lib/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifRole(roles ...string) gin.HandlerFunc {
    jwtToken := config.ReadENV()

    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            ctx.JSON(401, gin.H{"message": "Authorization header required"})
            ctx.Abort()
            return
        }

        tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

        token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
            return []byte(jwtToken), nil
        })

        if err != nil || !token.Valid {
            ctx.JSON(401, gin.H{"message": "Invalid token"})
            ctx.Abort()
            return
        }

        claims, _ := token.Claims.(jwt.MapClaims)
        userRole := claims["role"].(string)

        allowed := false
        for _, role := range roles {
            if userRole == role {
                allowed = true
                break
            }
        }

        if !allowed {
            ctx.JSON(403, gin.H{
                "success": false,
                "message": "Access role denied",
            })
            ctx.Abort()
            return
        }

        ctx.Next()
    }
}

