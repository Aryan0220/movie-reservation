package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/pashagolub/pgxmock/v2"
)

func TestReserveMovie_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/reserve", ReserveMovie)

	req := httptest.NewRequest("POST", "/reserve", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestReserveMovie_InvalidInput(t *testing.T) {
	app := fiber.New()
	app.Post("/reserve", func(c *fiber.Ctx) error {
		c.Locals("user_id", 1)
		return ReserveMovie(c)
	})

	req := httptest.NewRequest("POST", "/reserve", bytes.NewBufferString(`{"timetable_id":1,"screen_id":1,"seats":[],"date_time":""}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestReserveMovie_SeatUnavailable(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/reserve", func(c *fiber.Ctx) error {
		c.Locals("user_id", 1)
		return ReserveMovie(c)
	})

	dateTime := time.Date(2026, 6, 1, 15, 4, 0, 0, time.Local)
	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	schedule := []models.Schedule{{ScreenID: 2, Timings: []string{"15:04"}}}

	mock.ExpectQuery("SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes WHERE id=").
		WithArgs(10).
		WillReturnRows(pgxmock.NewRows([]string{"id", "movie_id", "schedule", "show_date", "normal_price", "vip_price"}).
			AddRow(10, 1, schedule, showDate, 10, 20))

	mock.ExpectQuery("SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=").
		WithArgs(2).
		WillReturnRows(pgxmock.NewRows([]string{"id", "auditorium_number", "normal_seats", "vip_seats", "type"}).
			AddRow(2, 1, []string{"A1"}, []string{"V1"}, "standard"))

	seatStatus := map[string]string{"A1": "booked"}
	mock.ExpectQuery("SELECT seat_status FROM show_seats WHERE showtime_id=").
		WithArgs(10, 2, "15:04").
		WillReturnRows(pgxmock.NewRows([]string{"seat_status"}).AddRow(seatStatus))

	payload := map[string]interface{}{
		"timetable_id": 10,
		"screen_id":    2,
		"seats":        []string{"A1"},
		"date_time":    dateTime.Format(time.RFC3339),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/reserve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusConflict {
		t.Fatalf("expected status 409, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCancelReservation_InvalidInput(t *testing.T) {
	app := fiber.New()
	app.Delete("/cancel", func(c *fiber.Ctx) error {
		c.Locals("user_id", 1)
		return CancelReservation(c)
	})

	req := httptest.NewRequest("DELETE", "/cancel", bytes.NewBufferString(`{"booking_id":0}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestCancelReservation_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Delete("/cancel", func(c *fiber.Ctx) error {
		c.Locals("user_id", 1)
		return CancelReservation(c)
	})

	dateTime := time.Date(2026, 6, 1, 15, 4, 0, 0, time.Local)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT user_id, showtime_id, screen_id, seats, booking_time FROM bookings WHERE id=").
		WithArgs(5).
		WillReturnRows(pgxmock.NewRows([]string{"user_id", "showtime_id", "screen_id", "seats", "booking_time"}).
			AddRow(1, 10, 2, []string{"A1"}, dateTime))
	mock.ExpectExec("UPDATE show_seats SET seat_status").
		WithArgs(10, 2, "15:04", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectExec("DELETE FROM bookings WHERE id=").
		WithArgs(5).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", "/cancel", bytes.NewBufferString(`{"booking_id":5}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetCapacity_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/capacity", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return GetCapacity(c)
	})

	req := httptest.NewRequest("POST", "/capacity", bytes.NewBufferString(`{"timetable_id":1,"screen_id":1,"date_time":"2026-06-01T15:04:00Z"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Fatalf("expected status 404, got %d", resp.StatusCode)
	}
}

func TestGetCapacity_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/capacity", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return GetCapacity(c)
	})

	ist, _ := time.LoadLocation("Asia/Kolkata")

	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, ist)
	dateTimeLocal := time.Date(2026, 6, 1, 15, 4, 0, 0, ist)
	payloadDateTime := dateTimeLocal.Format(time.RFC3339)
	dateTime, err := time.Parse(time.RFC3339, payloadDateTime)
	if err != nil {
		t.Fatalf("failed to parse date_time: %v", err)
	}
	schedule := []models.Schedule{{ScreenID: 2, Timings: []string{"15:04"}}}

	mock.ExpectQuery("SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes WHERE id=").
		WithArgs(10).
		WillReturnRows(pgxmock.NewRows([]string{"id", "movie_id", "schedule", "show_date", "normal_price", "vip_price"}).
			AddRow(10, 1, schedule, showDate, 10, 20))
	mock.ExpectQuery("SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=").
		WithArgs(2).
		WillReturnRows(pgxmock.NewRows([]string{"id", "auditorium_number", "normal_seats", "vip_seats", "type"}).
			AddRow(2, 1, []string{"A1", "A2"}, []string{"V1"}, "standard"))
	mock.ExpectQuery("SELECT seats FROM bookings WHERE showtime_id=").
		WithArgs(10, 2, dateTime).
		WillReturnRows(pgxmock.NewRows([]string{"seats"}).AddRow([]string{"A1"}))

	payload := map[string]interface{}{
		"timetable_id": 10,
		"screen_id":    2,
		"date_time":    payloadDateTime,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/capacity", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAllReservations_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Get("/bookings", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return GetAllReservations(c)
	})

	req := httptest.NewRequest("GET", "/bookings", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetAllReservations_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Get("/bookings", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return GetAllReservations(c)
	})

	ist, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, ist)
	rows := pgxmock.NewRows([]string{"id", "user_id", "showtime_id", "screen_id", "seats", "booking_time"}).
		AddRow(1, 2, 3, 4, []string{"A1"}, now)
	mock.ExpectQuery("SELECT id, user_id, showtime_id, screen_id, seats, booking_time FROM bookings ORDER BY booking_time DESC").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/bookings", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetRevenue_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/revenue", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return GetRevenue(c)
	})

	req := httptest.NewRequest("POST", "/revenue", bytes.NewBufferString(`{"movie_id":1}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != 500 {
		t.Fatalf("expected status 500, got %d", resp.StatusCode)
	}
}

func TestGetRevenue_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/revenue", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return GetRevenue(c)
	})

	rows := pgxmock.NewRows([]string{"screen_id", "seats", "normal_price", "vip_price"}).
		AddRow(2, []string{"A1", "V1"}, 10, 20)
	mock.ExpectQuery("SELECT b.screen_id, b.seats, mt.normal_price, mt.vip_price FROM bookings b JOIN showtimes mt").
		WithArgs(5).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=").
		WithArgs(2).
		WillReturnRows(pgxmock.NewRows([]string{"id", "auditorium_number", "normal_seats", "vip_seats", "type"}).
			AddRow(2, 1, []string{"A1"}, []string{"V1"}, "standard"))

	req := httptest.NewRequest("POST", "/revenue", bytes.NewBufferString(`{"movie_id":5}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestParseDateTime(t *testing.T) {
	value := "2026-06-01T15:04:00Z"
	parsed, err := parseDateTime(value)
	if err != nil {
		t.Fatalf("parseDateTime error: %v", err)
	}
	if parsed.IsZero() {
		t.Fatalf("expected parsed time")
	}
}
