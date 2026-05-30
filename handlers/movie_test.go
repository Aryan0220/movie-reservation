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

func TestAddMovie_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Post("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return AddMovie(c)
	})

	req := httptest.NewRequest("POST", "/movie", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestAddMovie_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return AddMovie(c)
	})

	movie := models.Movie{Title: "Dune", Description: "Sci-fi", PosterURL: "url", Genre: []string{"Sci-Fi"}}
	payload, _ := json.Marshal(movie)

	mock.ExpectExec("INSERT INTO movies").
		WithArgs(movie.Title, movie.Description, movie.PosterURL, movie.Genre).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	req := httptest.NewRequest("POST", "/movie", bytes.NewReader(payload))
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

func TestUpdateMovie_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Patch("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return UpdateMovie(c)
	})

	req := httptest.NewRequest("PATCH", "/movie", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestUpdateMovie_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Patch("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return UpdateMovie(c)
	})

	movie := models.Movie{Title: "Dune", Description: "Updated", PosterURL: "url2", Genre: []string{"Sci-Fi"}}
	payload, _ := json.Marshal(movie)

	mock.ExpectExec("UPDATE movies").
		WithArgs(movie.Title, movie.Description, movie.PosterURL, movie.Genre, movie.Title).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := httptest.NewRequest("PATCH", "/movie", bytes.NewReader(payload))
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

func TestDeleteMovie_Unauthorized(t *testing.T) {
	app := fiber.New()
	app.Delete("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", false)
		return DeleteMovie(c)
	})

	req := httptest.NewRequest("DELETE", "/movie", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403, got %d", resp.StatusCode)
	}
}

func TestDeleteMovie_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Delete("/movie", func(c *fiber.Ctx) error {
		c.Locals("Role", true)
		return DeleteMovie(c)
	})

	payload := []byte(`{"title":"Dune"}`)
	mock.ExpectExec("DELETE FROM movies").
		WithArgs("Dune").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	req := httptest.NewRequest("DELETE", "/movie", bytes.NewReader(payload))
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

func TestGetMovies_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Get("/movies", GetMovies)

	rows := pgxmock.NewRows([]string{"title"}).AddRow("Dune")
	mock.ExpectQuery("SELECT title FROM MOVIES").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/movies", nil)
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

func TestViewSeats_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/seats", ViewSeats)

	payload := []byte(`{"showtime_id":1,"screen_id":2,"show_time":"15:30"}`)
	status := map[string]string{"A1": "available"}
	rows := pgxmock.NewRows([]string{"seat_status"}).AddRow(status)
	mock.ExpectQuery("SELECT seat_status FROM show_seats").
		WithArgs(1, 2, "15:30").
		WillReturnRows(rows)

	req := httptest.NewRequest("POST", "/seats", bytes.NewReader(payload))
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

func TestGetMovieTimings_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/timings", GetMovieTimings)

	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
	schedule := []models.Schedule{{ScreenID: 1, Timings: []string{"15:00"}}}
	rows := pgxmock.NewRows([]string{"id", "movie_id", "schedule", "show_date", "normal_price", "vip_price"}).
		AddRow(1, 10, schedule, showDate, 10, 20)
	mock.ExpectQuery("SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes").
		WithArgs("Dune").
		WillReturnRows(rows)

	req := httptest.NewRequest("POST", "/timings", bytes.NewBufferString(`"Dune"`))
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
