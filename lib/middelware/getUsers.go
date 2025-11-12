package middelware

import (
	"back-end-coffeShop/lib/config"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserFromToken(ctx *gin.Context) int {
	jwtToken := config.ReadENV()

	authHeader := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &LoginClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return []byte(jwtToken), nil
	})

	if err != nil || !token.Valid {
		return 0
	}


	return claims.UserID
}
