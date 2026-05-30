package services

import (
	"reflect"
	"testing"
	"time"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/pashagolub/pgxmock/v2"
)

func TestCreateUser(t *testing.T) {
	mock := testutils.NewMockDB(t)
	user := models.User{Name: "A", Email: "a@example.com", Role: false, Password: "secret"}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.Role, pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	if err := CreateUser(user); err != nil {
		t.Fatalf("CreateUser error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetUserByEmail(t *testing.T) {
	mock := testutils.NewMockDB(t)

	rows := pgxmock.NewRows([]string{"id", "name", "email", "admin", "password"}).
		AddRow(1, "A", "a@example.com", true, "hash")
	mock.ExpectQuery("SELECT id, name, email, admin, password FROM users WHERE email=").
		WithArgs("a@example.com").
		WillReturnRows(rows)

	user, err := GetUserByEmail("a@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail error: %v", err)
	}
	if user.ID != 1 || user.Email != "a@example.com" || user.Role != true {
		t.Fatalf("unexpected user: %+v", user)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetMovieTimings(t *testing.T) {
	mock := testutils.NewMockDB(t)
	showDate := time.Date(2026, 5, 30, 0, 0, 0, 0, time.Local)
	schedule := []models.Schedule{{ScreenID: 1, Timings: []string{"15:00"}}}

	rows := pgxmock.NewRows([]string{"id", "movie_id", "schedule", "show_date", "normal_price", "vip_price"}).
		AddRow(10, 5, schedule, showDate, 10, 20)
	mock.ExpectQuery("SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes").
		WithArgs("Dune").
		WillReturnRows(rows)

	timings, err := GetMovieTimings("Dune")
	if err != nil {
		t.Fatalf("GetMovieTimings error: %v", err)
	}
	if len(timings) != 1 || !reflect.DeepEqual(timings[0].Schedule, schedule) {
		t.Fatalf("unexpected timings: %+v", timings)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
