package handlers

import (
	"booking-system/models"
	"booking-system/services"
	"booking-system/config"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	config.PrintLog("Registering user: "+user.Email, "INFO")

	err := services.CreateUser(user)
	if err != nil {
		config.PrintLog("Error creating user: "+err.Error(), "ERROR")
		return c.Status(500).JSON(fiber.Map{"message": "User Creation Failed"})
	}

	config.PrintLog("User created successfully: "+user.Email, "INFO")
	return c.JSON(fiber.Map{"message": "User created"})
}

func Login(c *fiber.Ctx) error {
	var input models.User

	c.BodyParser(&input)

	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		config.PrintLog("Error fetching user: "+err.Error(), "ERROR")
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Credentials", "details": err.Error()})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		config.PrintLog("Invalid credentials for user: "+user.Email, "ERROR")
		return c.Status(401).JSON(fiber.Map{"error": "invalid Credentials"})
	}

	token, _ := services.GenerateToken(user.ID, user.Role)
	config.PrintLog("User logged in successfully: "+user.Email, "INFO")
	return c.JSON(fiber.Map{"token": token})
}

func Promote(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.PromoteToAdmin(input)
	if err != nil {
		config.PrintLog("Error promoting user: "+err.Error(), "ERROR")
		return c.Status(500).JSON(fiber.Map{"error": "Failed to Promote User"})
	}

	config.PrintLog("User promoted to admin: "+input.Email, "INFO")
	return c.Status(200).JSON(fiber.Map{"message": "User Promoted to Admin"})
}
