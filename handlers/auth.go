package handlers

import (
	"booking-system/models"
	"booking-system/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.CreateUser(user)
	if err != nil {
		log.Print("Error creating user: ", err)
		return c.Status(500).JSON(fiber.Map{"message": "User Creation Failed"})
	}

	return c.JSON(fiber.Map{"message": "User created"})
}

func Login(c *fiber.Ctx) error {
	var input models.User

	c.BodyParser(&input)

	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid Credentials"})
	}

	token, _ := services.GenerateToken(user.ID, user.Role)

	return c.JSON(fiber.Map{"token": token})
}

func Promote(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.PromoteToAdmin(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Promote User"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "User Promoted to Admin"})
}
