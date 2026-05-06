package middleware

import (
	"booking-system/config"
	"strings"
	"log"
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


	tokenStr := parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invail or expired token",
		})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		log.Println(claims, claims["role"])
		c.Locals("user_id", claims["user_id"])
		c.Locals("Role", claims["role"])
	}

	return c.Next()
}
