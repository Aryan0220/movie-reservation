package handlers

import (
	"booking-system/config"
	"booking-system/models"
	"booking-system/services"
	"time"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

func AddShowTime(c *fiber.Ctx) error {
	role := c.Locals("Role").(bool)
	if !role {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized", "details": role})
	}

	var input models.MovieTimetable
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input", "details": err.Error()})
	}

	if err := validateTimetableInput(input, false); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.AddShowTime(input); err != nil { 
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add showtime"})
	}
	config.PrintLog("Showtime added successfully for movie id "+strconv.Itoa(input.MovieID), "INFO")
	return c.JSON(fiber.Map{"message": "Showtime added successfully"})
}

func UpdateShowTime(c *fiber.Ctx) error {
	role := c.Locals("Role").(bool)
	if !role {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}

	var input models.MovieTimetable
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validateTimetableInput(input, true); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := services.UpdateShowTime(input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update showtime"})
	}
	config.PrintLog("Showtime updated successfully for timetable id "+strconv.Itoa(input.ID), "INFO")
	return c.JSON(fiber.Map{"message": "Showtime updated successfully"})
}

func validateTimetableInput(input models.MovieTimetable, requireID bool) error {
	if requireID && input.ID <= 0 {
		return fiber.NewError(400, "Missing timetable id")
	}
	if input.MovieID <= 0 {
		return fiber.NewError(400, "Missing movie id")
	}
	if len(input.Schedule) == 0 {
		return fiber.NewError(400, "At least one schedule entry is required")
	}
	if input.NormalPrice <= 0 || input.VipPrice <= 0 {
		return fiber.NewError(400, "Prices must be greater than zero")
	}

	showDate := input.ShowDate
	
	today := time.Now().In(time.Local)
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	if showDate.Before(todayDate) {
		return fiber.NewError(400, "Show date cannot be in the past")
	}

	seenScreens := make(map[int]struct{}, len(input.Schedule))
	for _, schedule := range input.Schedule {
		if schedule.ScreenID <= 0 {
			return fiber.NewError(400, "Screen ids must be positive")
		}
		if _, exists := seenScreens[schedule.ScreenID]; exists {
			return fiber.NewError(400, "Duplicate screen found")
		}
		seenScreens[schedule.ScreenID] = struct{}{}
	}

	config.PrintLog("Timetable input validated successfully for movie id "+strconv.Itoa(input.MovieID), "INFO")
	return nil
}
