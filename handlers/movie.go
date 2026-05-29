package handlers

import (
	"booking-system/models"
	"booking-system/services"
	"booking-system/config"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddMovie(c *fiber.Ctx) error {
	role := c.Locals("Role").(bool)
	if !role {
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

	config.PrintLog("Movie added successfully: "+input.Title, "INFO")
	return c.JSON(fiber.Map{"message": "Movie added successfully"})
}

func UpdateMovie(c *fiber.Ctx) error {
	role := c.Locals("Role").(bool)

	if !role {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}
	
	var input models.Movie

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.Update_Movie(input)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update movie"})
	}
	config.PrintLog("Movie updated successfully: "+input.Title, "INFO")
	return c.JSON(fiber.Map{"message": "Movie Updated"})
}

func DeleteMovie(c *fiber.Ctx) error {
	role := c.Locals("Role").(bool)

	if !role {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}
	
	var input models.Movie

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.Delete_Movie(input)
	
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete movie"})
	}
	config.PrintLog("Movie deleted successfully: "+input.Title, "INFO")
	return c.JSON(fiber.Map{"message": "Movie Deleted"})
}

func GetMovieTimings(c *fiber.Ctx) error {
	var input string

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	timings, err := services.GetMovieTimings(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get movie timings", "details": err.Error()})
	}
	config.PrintLog("Movie timings fetched for: "+input, "INFO")
	return c.Status(200).JSON(timings)
}

func ViewSeats(c *fiber.Ctx) error {
	var input services.SeatStatusRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	seat_status := services.ViewSeats(input)

	if seat_status == nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get seat status"})
	}
	config.PrintLog("Seat status fetched for timetable "+strconv.Itoa(input.ShowTimeID)+", screen "+strconv.Itoa(input.ScreenID), "INFO")
	return c.Status(200).JSON(seat_status)
}

func GetMovies(c *fiber.Ctx) error {
	movies, err := services.Get_Movies()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get movies"})
	}
	return c.Status(200).JSON(movies)
}