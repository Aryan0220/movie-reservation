package handlers

import (
	"booking-system/models"
	"booking-system/services"
	"github.com/gofiber/fiber/v2"
)

func AddMovie(c *fiber.Ctx) error {
	role := c.Locals("Role")
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}

	var input models.Movie

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.Add_Movie(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add movie"})
	}

	return c.JSON(fiber.Map{"message": "Movie added successfully"})
}

func UpdateMovie(c *fiber.Ctx) error {
	role := c.Locals("Role")

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}
	
	var input models.Movie

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.Update_Movie(input)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add movie"})
	}

	return c.JSON(fiber.Map{"message": "Movie Updated"})
}

func DeleteMovie(c *fiber.Ctx) error {
	role := c.Locals("Role")

	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}
	
	var input models.Movie

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.Delete_Movie(input)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add movie"})
	}

	return c.JSON(fiber.Map{"message": "Movie Deleted"})
}
