package middleware

import (
	"booking-system/config"
	"encoding/json"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Protected(c *fiber.Ctx) error {
	header := c.Get("Authorization")

	if header == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	parts := strings.Split(header, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization format",
		})
	}


	secret := config.GetEnv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server misconfiguration",
		})
	}

	tokenStr := parts[1]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	role, _ := claims["role"].(bool)

	userID, ok := toIntClaim(claims["user_id"])
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	c.Locals("user_id", userID)
	c.Locals("Role", role)

	return c.Next()
}

func toIntClaim(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case int64:
		return int(v), true
	case float64:
		return int(v), true
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return 0, false
		}
		return int(parsed), true
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}
