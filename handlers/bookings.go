package handlers

import (
	"booking-system/services"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type reservationRequest struct {
	TimetableID int      `json:"timetable_id"`
	ScreenID    int      `json:"screen_id"`
	Seats       []string `json:"seats"`
	DateTime    string   `json:"date_time"`
}

type cancelRequest struct {
	BookingID int `json:"booking_id"`
}

type capacityRequest struct {
	TimetableID int    `json:"timetable_id"`
	ScreenID    int    `json:"screen_id"`
	DateTime    string `json:"date_time"`
}

type revenueRequest struct {
	MovieID int `json:"movie_id"`
}

func ReserveMovie(c *fiber.Ctx) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var input reservationRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if input.TimetableID <= 0 || input.ScreenID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing timetable or screen id"})
	}
	if len(input.Seats) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "At least one seat is required"})
	}

	dateTime, err := parseDateTime(input.DateTime)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date_time", "details": err.Error()})
	}

	err = services.ReserveTicket(userID, input.TimetableID, input.ScreenID, input.Seats, dateTime, time.Now())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrSeatUnavailable):
			return c.Status(409).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, services.ErrInvalidShowtime), errors.Is(err, services.ErrPastReservation):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, services.ErrNotFound):
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "Failed to reserve seat"})
		}
	}

	return c.JSON(fiber.Map{"message": "Reservation created"})
}

func CancelReservation(c *fiber.Ctx) error {
	userID, err := requireUserID(c)
	if err != nil {
		return err
	}

	var input cancelRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.BookingID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing booking id"})
	}

	err = services.CancelReservation(userID, input.BookingID, time.Now())
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPastReservation):
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, services.ErrNotOwner):
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		case errors.Is(err, services.ErrNotFound):
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel reservation"})
		}
	}

	return c.JSON(fiber.Map{"message": "Reservation canceled"})
}

func GetCapacity(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	var input capacityRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.TimetableID <= 0 || input.ScreenID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing timetable or screen id"})
	}

	dateTime, err := parseDateTime(input.DateTime)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date_time", "details": err.Error()})
	}

	total, available, err := services.GetCapacity(input.TimetableID, input.ScreenID, dateTime)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound), errors.Is(err, services.ErrInvalidShowtime):
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "Failed to get capacity"})
		}
	}

	return c.JSON(fiber.Map{"total": total, "available": available})
}

func GetAllReservations(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	bookings, err := services.GetAllBookings()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get bookings"})
	}

	return c.JSON(bookings)
}

func GetRevenue(c *fiber.Ctx) error {
	if err := requireAdmin(c); err != nil {
		return err
	}

	var input revenueRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if input.MovieID <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Missing movie id"})
	}

	revenue, err := services.GetMovieRevenue(input.MovieID, time.Now())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get revenue"})
	}

	return c.JSON(fiber.Map{"movie_id": input.MovieID, "revenue": revenue})
}

func requireAdmin(c *fiber.Ctx) error {
	role := c.Locals("Role")
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "User not authorized"})
	}
	return nil
}

func requireUserID(c *fiber.Ctx) (int, error) {
	value := c.Locals("user_id")
	if value == nil {
		return 0, c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return parsed, nil
	default:
		return 0, c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
}

func parseDateTime(input string) (time.Time, error) {
	if input == "" {
		return time.Time{}, errors.New("date_time is required")
	}

	if parsed, err := time.Parse(time.RFC3339, input); err == nil {
		return parsed, nil
	}

	layout := "2006-01-02 15:04:05"
	return time.ParseInLocation(layout, input, time.Local)
}
