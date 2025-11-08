package middelware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginClaims struct{
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(JWTtoken string, role string) (string, error){
	claims := LoginClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTtoken))
	
	return tokenString, err
}