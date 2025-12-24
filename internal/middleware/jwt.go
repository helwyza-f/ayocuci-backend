// internal/middleware/jwt.go
package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, role string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
