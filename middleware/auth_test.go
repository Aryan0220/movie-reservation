package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"booking-system/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func setupProtectedApp() *fiber.App {
	app := fiber.New()
	app.Get("/protected", Protected, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user_id": c.Locals("user_id"),
			"role":    c.Locals("Role"),
		})
	})
	return app
}

func TestProtected_MissingHeader(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/protected", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestProtected_InvalidHeaderFormat(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Token abc")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestProtected_InvalidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-token")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestProtected_InvalidAlg(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	claims := jwt.MapClaims{
		"user_id": 123,
		"role":    "admin",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestProtected_ValidToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	token, err := services.GenerateToken(123, "admin")
	if err != nil {
		t.Fatalf("GenerateToken error: %v", err)
	}

	app := setupProtectedApp()
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["role"] != "admin" {
		t.Fatalf("expected role admin, got %v", body["role"])
	}
	if body["user_id"] != float64(123) {
		t.Fatalf("expected user_id 123, got %v", body["user_id"])
	}
}
