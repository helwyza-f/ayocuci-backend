package common

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(
	userID uint,
	clientID uint,
	role string,
	secret string,
) (string, error) {

	claims := jwt.MapClaims{
		"user_id":   userID,
		"client_id": clientID,
		"role":      role,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
