package services

import (
	"time"
	"booking-system/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, userRole string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role": userRole,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
}
