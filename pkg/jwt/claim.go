package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func generateClaims(userID int, expires time.Duration) *Claims {
	return &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}
