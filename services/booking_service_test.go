package services

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"booking-system/models"
	"booking-system/testutils"

	"github.com/pashagolub/pgxmock/v2"
)

func TestNormalizeSeats(t *testing.T) {
	result, err := normalizeSeats([]string{"A1", "A2"})
	if err != nil {
		t.Fatalf("normalizeSeats error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("unexpected seats: %+v", result)
	}

	if _, err := normalizeSeats([]string{""}); err == nil {
		t.Fatalf("expected error for empty seat")
	}
	if _, err := normalizeSeats([]string{"A1", "A1"}); err == nil {
		t.Fatalf("expected error for duplicate seat")
	}
}

func TestValidateSeatsExist(t *testing.T) {
	seats := map[string]struct{}{"A1": {}, "V1": {}}
	if err := validateSeatsExist(seats, []string{"A1"}, []string{"V1"}); err != nil {
		t.Fatalf("validateSeatsExist error: %v", err)
	}
	if err := validateSeatsExist(map[string]struct{}{"X": {}}, []string{"A1"}, []string{"V1"}); err == nil {
		t.Fatalf("expected error for unknown seat")
	}
}

func TestMatchesTimingAndDate(t *testing.T) {
	date := time.Date(2026, 6, 1, 15, 4, 0, 0, time.Local)
	if !matchesTiming(date, []string{"15:04"}) {
		t.Fatalf("expected timing match")
	}
	if matchesTiming(date, []string{"16:00"}) {
		t.Fatalf("expected timing mismatch")
	}
	if !isSameDate(date, time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)) {
		t.Fatalf("expected same date")
	}
}

func TestBuildSeatStatusPatch(t *testing.T) {
	patch, err := buildSeatStatusPatch(map[string]struct{}{"A1": {}, "A2": {}}, "booked")
	if err != nil {
		t.Fatalf("buildSeatStatusPatch error: %v", err)
	}
	if patch == "" {
		t.Fatalf("expected patch to be non-empty")
	}
}

func TestReserveTicket_SeatUnavailable(t *testing.T) {
	mock := testutils.NewMockDB(t)
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
			AddRow(2, 1, []string{"A1", "A2"}, []string{"V1"}, "standard"))

	seatStatus := map[string]string{"A1": "booked", "A2": "available"}
	mock.ExpectQuery("SELECT seat_status FROM show_seats WHERE showtime_id=").
		WithArgs(10, 2, "15:04").
		WillReturnRows(pgxmock.NewRows([]string{"seat_status"}).AddRow(seatStatus))

	err := ReserveTicket(1, 10, 2, []string{"A1"}, dateTime, dateTime.Add(-time.Hour))
	if !errors.Is(err, ErrSeatUnavailable) {
		t.Fatalf("expected ErrSeatUnavailable, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestReserveTicket_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
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
			AddRow(2, 1, []string{"A1", "A2"}, []string{"V1"}, "standard"))

	seatStatus := map[string]string{"A1": "available", "A2": "available"}
	mock.ExpectQuery("SELECT seat_status FROM show_seats WHERE showtime_id=").
		WithArgs(10, 2, "15:04").
		WillReturnRows(pgxmock.NewRows([]string{"seat_status"}).AddRow(seatStatus))

	mock.ExpectExec("INSERT INTO bookings").
		WithArgs(1, 10, 2, []string{"A1"}, dateTime).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("UPDATE show_seats SET seat_status").
		WithArgs(10, 2, "15:04", pgxmock.AnyArg()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := ReserveTicket(1, 10, 2, []string{"A1"}, dateTime, dateTime.Add(-time.Hour))
	if err != nil {
		t.Fatalf("ReserveTicket error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestCancelReservation_Success(t *testing.T) {
	mock := testutils.NewMockDB(t)
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

	err := CancelReservation(1, 5, dateTime.Add(-time.Hour))
	if err != nil {
		t.Fatalf("CancelReservation error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetCapacity(t *testing.T) {
	mock := testutils.NewMockDB(t)
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
			AddRow(2, 1, []string{"A1", "A2"}, []string{"V1"}, "standard"))

	rows := pgxmock.NewRows([]string{"seats"}).AddRow([]string{"A1"})
	mock.ExpectQuery("SELECT seats FROM bookings WHERE showtime_id=").
		WithArgs(10, 2, dateTime).
		WillReturnRows(rows)

	total, available, err := GetCapacity(10, 2, dateTime)
	if err != nil {
		t.Fatalf("GetCapacity error: %v", err)
	}
	if total != 3 || available != 2 {
		t.Fatalf("unexpected capacity: total=%d available=%d", total, available)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetAllBookings(t *testing.T) {
	mock := testutils.NewMockDB(t)
	now := time.Date(2026, 6, 1, 0, 0, 0, 0, time.Local)

	rows := pgxmock.NewRows([]string{"id", "user_id", "showtime_id", "screen_id", "seats", "booking_time"}).
		AddRow(1, 2, 3, 4, []string{"A1"}, now)
	mock.ExpectQuery("SELECT id, user_id, showtime_id, screen_id, seats, booking_time FROM bookings").
		WillReturnRows(rows)

	bookings, err := GetAllBookings()
	if err != nil {
		t.Fatalf("GetAllBookings error: %v", err)
	}
	if len(bookings) != 1 || bookings[0].UserID != 2 {
		t.Fatalf("unexpected bookings: %+v", bookings)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetMovieRevenue(t *testing.T) {
	mock := testutils.NewMockDB(t)

	rows := pgxmock.NewRows([]string{"screen_id", "seats", "normal_price", "vip_price"}).
		AddRow(2, []string{"A1", "V1"}, 10, 20)
	mock.ExpectQuery("SELECT b.screen_id, b.seats, mt.normal_price, mt.vip_price FROM bookings b JOIN showtimes mt").
		WithArgs(5).
		WillReturnRows(rows)

	mock.ExpectQuery("SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=").
		WithArgs(2).
		WillReturnRows(pgxmock.NewRows([]string{"id", "auditorium_number", "normal_seats", "vip_seats", "type"}).
			AddRow(2, 1, []string{"A1"}, []string{"V1"}, "standard"))

	revenue, err := GetMovieRevenue(5, time.Now())
	if err != nil {
		t.Fatalf("GetMovieRevenue error: %v", err)
	}
	if revenue != 30 {
		t.Fatalf("unexpected revenue: %d", revenue)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestGetTimingsByScreenID(t *testing.T) {
	data := []models.Schedule{{ScreenID: 1, Timings: []string{"10:00"}}}
	if timings := getTimingsByScreenID(1, data); !reflect.DeepEqual(timings, []string{"10:00"}) {
		t.Fatalf("unexpected timings: %+v", timings)
	}
	if timings := getTimingsByScreenID(2, data); timings != nil {
		t.Fatalf("expected nil timings")
	}
}
