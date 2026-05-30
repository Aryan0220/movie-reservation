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

func TestAddShowTime_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/timetable", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return AddShowTime(c)
	})

	req := httptest.NewRequest("POST", "/timetable", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestAddShowTime_ValidationError(t *testing.T) {
	app := fiber.New()
	app.Post("/timetable", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return AddShowTime(c)
	})

	payload := []byte(`{"movie_id":0}`)
	req := httptest.NewRequest("POST", "/timetable", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestAddShowTime_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/timetable", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return AddShowTime(c)
	})
	ist, _ := time.LoadLocation("Asia/Kolkata")

	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, ist)
	payloadShowDate := showDate.Format(time.RFC3339)
	input := models.MovieTimetable{
		MovieID:     10,
		ShowDate:    showDate,
		NormalPrice: 10,
		VipPrice:    20,
		Schedule: []models.Schedule{{
			ScreenID: 2,
			Timings:  []string{"15:04", "18:00"},
		}},
	}
	payload, _ := json.Marshal(input)
	expectedShowDate, err := time.Parse(time.RFC3339, payloadShowDate)
	if err != nil {
		t.Fatalf("failed to parse show_date: %v", err)
	}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO showtimes").
		WithArgs(input.MovieID, input.Schedule, expectedShowDate, input.NormalPrice, input.VipPrice).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(99))

	mock.ExpectQuery("SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=").
		WithArgs(2).
		WillReturnRows(pgxmock.NewRows([]string{"id", "auditorium_number", "normal_seats", "vip_seats", "type"}).
			AddRow(2, 1, []string{"A1", "A2"}, []string{"V1"}, "standard"))

	mock.ExpectExec("INSERT INTO show_seats").
		WithArgs(99, 2, "15:04", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO show_seats").
		WithArgs(99, 2, "18:00", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	req := httptest.NewRequest("POST", "/timetable", bytes.NewReader(payload))
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

func TestUpdateShowTime_ValidationError(t *testing.T) {
	app := fiber.New()
	app.Patch("/timetable", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return UpdateShowTime(c)
	})

	payload := []byte(`{"id":0,"movie_id":1}`)
	req := httptest.NewRequest("PATCH", "/timetable", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestUpdateShowTime_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Patch("/timetable", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return UpdateShowTime(c)
	})

	ist, _ := time.LoadLocation("Asia/Kolkata")

	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, ist)
	payloadShowDate := showDate.Format(time.RFC3339)
	input := models.MovieTimetable{ID: 5, MovieID: 10, ShowDate: showDate, NormalPrice: 10, VipPrice: 20, Schedule: []models.Schedule{{ScreenID: 1, Timings: []string{"12:00"}}}}
	payload, _ := json.Marshal(input)
	expectedShowDate, err := time.Parse(time.RFC3339, payloadShowDate)
	if err != nil {
		t.Fatalf("failed to parse show_date: %v", err)
	}
	mock.ExpectExec("UPDATE showtimes SET").
		WithArgs(input.MovieID, input.Schedule, expectedShowDate, input.NormalPrice, input.VipPrice, input.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := httptest.NewRequest("PATCH", "/timetable", bytes.NewReader(payload))
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
