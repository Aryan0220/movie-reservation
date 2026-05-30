package services

import (
	"reflect"
	"testing"
	"time"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/pashagolub/pgxmock/v2"
)

func TestBuildSeatStatus(t *testing.T) {
	status := buildSeatStatus([]string{"A1", "A2"}, []string{"V1"})
	want := map[string]string{"A1": "available", "A2": "available", "V1": "available"}
	if !reflect.DeepEqual(status, want) {
		t.Fatalf("unexpected seat status: %+v", status)
	}
}

func TestNormalizeShowTime(t *testing.T) {
	cases := map[string]string{
		"15:04:05": "15:04",
		"15:04":    "15:04",
		"bad":      "bad",
	}
	for input, want := range cases {
		if got := normalizeShowTime(input); got != want {
			t.Fatalf("normalizeShowTime(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestAddShowTime(t *testing.T) {
	mock := testutils.NewMockDB(t)
	showDate := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)
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

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO showtimes").
		WithArgs(input.MovieID, input.Schedule, input.ShowDate, input.NormalPrice, input.VipPrice).
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

	if err := AddShowTime(input); err != nil {
		t.Fatalf("AddShowTime error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUpdateShowTime(t *testing.T) {
	mock := testutils.NewMockDB(t)
	input := models.MovieTimetable{ID: 5, MovieID: 10, NormalPrice: 10, VipPrice: 20}

	mock.ExpectExec("UPDATE showtimes SET").
		WithArgs(input.MovieID, input.Schedule, input.ShowDate, input.NormalPrice, input.VipPrice, input.ID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	if err := UpdateShowTime(input); err != nil {
		t.Fatalf("UpdateShowTime error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
