package services

import (
	"time"
	"booking-system/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, userRole bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role": userRole,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
}
