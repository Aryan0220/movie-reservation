package handlers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/pashagolub/pgxmock/v2"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister_InvalidInput(t *testing.T) {
	app := fiber.New()
	app.Post("/register", Register)

	req := httptest.NewRequest("POST", "/register", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestRegister_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/register", Register)

	user := models.User{Name: "A", Email: "a@example.com", Password: "secret", Role: false}
	payload, _ := json.Marshal(user)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.Role, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)
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

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Post("/login", Login)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), 4)
	rows := pgxmock.NewRows([]string{"id", "name", "email", "admin", "password"}).
		AddRow(1, "A", "a@example.com", false, string(hash))
	mock.ExpectQuery("SELECT id, name, email, admin, password FROM users WHERE email=").
		WithArgs("a@example.com").
		WillReturnRows(rows)

	payload := []byte(`{"email":"a@example.com","password":"wrong"}`)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	t.Setenv("JWT_SECRET", "test-secret")
	app := fiber.New()
	app.Post("/login", Login)

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), 4)
	rows := pgxmock.NewRows([]string{"id", "name", "email", "admin", "password"}).
		AddRow(1, "A", "a@example.com", true, string(hash))
	mock.ExpectQuery("SELECT id, name, email, admin, password FROM users WHERE email=").
		WithArgs("a@example.com").
		WillReturnRows(rows)

	payload := []byte(`{"email":"a@example.com","password":"secret"}`)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(payload))
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

func TestPromote_InvalidInput(t *testing.T) {
	app := fiber.New()
	app.Patch("/promote", Promote)

	req := httptest.NewRequest("PATCH", "/promote", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", resp.StatusCode)
	}
}

func TestPromote_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
	app := fiber.New()
	app.Patch("/promote", Promote)

	payload := []byte(`{"email":"a@example.com"}`)
	mock.ExpectExec("UPDATE users SET admin=true").
		WithArgs("a@example.com").
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	req := httptest.NewRequest("PATCH", "/promote", bytes.NewReader(payload))
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
