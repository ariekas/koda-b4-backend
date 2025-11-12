package middelware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LoginClaims struct{
	UserID int    `json:"user_id"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(JWTtoken string, role string,  userID int) (string, error){
	claims := LoginClaims{
		UserID: userID,
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