package services

import (
	"booking-system/config"
	"booking-system/models"
	"context"
	"errors"
	"time"
)

var (
	ErrSeatUnavailable = errors.New("seat unavailable")
	ErrInvalidShowtime = errors.New("invalid showtime")
	ErrNotFound        = errors.New("record not found")
	ErrPastReservation = errors.New("reservation is in the past")
	ErrNotOwner        = errors.New("not booking owner")
)

func ReserveTicket(userID, timetableID, screenID int, seats []string, dateTime time.Time, now time.Time) error {
	if dateTime.Before(now) {
		return ErrPastReservation
	}

	timetable, err := getTimetableByID(timetableID)
	if err != nil {
		return err
	}

	if !isSameDate(timetable.ShowDate, dateTime) || !matchesTiming(dateTime, timetable.Timings) {
		return ErrInvalidShowtime
	}

	if !containsInt(timetable.Screens, screenID) {
		return ErrInvalidShowtime
	}

	screen, err := getScreenByID(screenID)
	if err != nil {
		return err
	}

	requestedSeats, err := normalizeSeats(seats)
	if err != nil {
		return err
	}

	if err := validateSeatsExist(requestedSeats, screen.Normal, screen.Vip); err != nil {
		return err
	}

	bookedSeats, err := getBookedSeats(timetableID, screenID, dateTime)
	if err != nil {
		return err
	}

	for seat := range requestedSeats {
		if _, exists := bookedSeats[seat]; exists {
			return ErrSeatUnavailable
		}
	}

	_, err = config.DB.Exec(context.Background(),
		"INSERT INTO bookings (user_id, timetable_id, screen_id, reservation, date_and_time) VALUES ($1, $2, $3, $4, $5)",
		userID, timetableID, screenID, seats, dateTime,
	)
	return err
}

func CancelReservation(userID, bookingID int, now time.Time) error {
	var ownerID int
	var dateTime time.Time

	err := config.DB.QueryRow(context.Background(),
		"SELECT user_id, date_and_time FROM bookings WHERE id=$1",
		bookingID,
	).Scan(&ownerID, &dateTime)
	if err != nil {
		return ErrNotFound
	}

	if ownerID != userID {
		return ErrNotOwner
	}

	if !dateTime.After(now) {
		return ErrPastReservation
	}

	_, err = config.DB.Exec(context.Background(),
		"DELETE FROM bookings WHERE id=$1",
		bookingID,
	)
	return err
}

func GetCapacity(timetableID, screenID int, dateTime time.Time) (int, int, error) {
	timetable, err := getTimetableByID(timetableID)
	if err != nil {
		return 0, 0, err
	}

	if !isSameDate(timetable.ShowDate, dateTime) || !matchesTiming(dateTime, timetable.Timings) {
		return 0, 0, ErrInvalidShowtime
	}

	if !containsInt(timetable.Screens, screenID) {
		return 0, 0, ErrInvalidShowtime
	}

	screen, err := getScreenByID(screenID)
	if err != nil {
		return 0, 0, err
	}

	bookedSeats, err := getBookedSeats(timetableID, screenID, dateTime)
	if err != nil {
		return 0, 0, err
	}

	total := len(screen.Normal) + len(screen.Vip)
	available := total - len(bookedSeats)
	if available < 0 {
		available = 0
	}

	return total, available, nil
}

func GetAllBookings() ([]models.BookingDetail, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT id, user_id, timetable_id, screen_id, reservation, date_and_time FROM bookings ORDER BY date_and_time DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.BookingDetail
	for rows.Next() {
		var booking models.BookingDetail
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.TimetableID, &booking.ScreenID, &booking.Reservation, &booking.DateTime); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func GetMovieRevenue(movieID int, now time.Time) (int, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT b.screen_id, b.reservation, mt.normal_price, mt.vip_price FROM bookings b JOIN movie_timetables mt ON mt.id = b.timetable_id WHERE mt.movie_id=$1 AND b.date_and_time <= $2",
		movieID, now,
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	seatCache := make(map[int]models.Screen)
	revenue := 0

	for rows.Next() {
		var screenID int
		var seats []string
		var normalPrice int
		var vipPrice int
		if err := rows.Scan(&screenID, &seats, &normalPrice, &vipPrice); err != nil {
			return 0, err
		}

		screen, ok := seatCache[screenID]
		if !ok {
			loaded, err := getScreenByID(screenID)
			if err != nil {
				return 0, err
			}
			screen = loaded
			seatCache[screenID] = loaded
		}

		normalSet := toStringSet(screen.Normal)
		vipSet := toStringSet(screen.Vip)

		normalCount := 0
		vipCount := 0
		for _, seat := range seats {
			if _, exists := normalSet[seat]; exists {
				normalCount++
				continue
			}
			if _, exists := vipSet[seat]; exists {
				vipCount++
			}
		}

		revenue += normalCount*normalPrice + vipCount*vipPrice
	}

	return revenue, nil
}

func getTimetableByID(timetableID int) (models.MovieTimetable, error) {
	var timetable models.MovieTimetable
	err := config.DB.QueryRow(context.Background(),
		"SELECT id, movie_id, timings, screens, show_date, normal_price, vip_price FROM movie_timetables WHERE id=$1",
		timetableID,
	).Scan(&timetable.ID, &timetable.MovieID, &timetable.Timings, &timetable.Screens, &timetable.ShowDate, &timetable.NormalPrice, &timetable.VipPrice)
	if err != nil {
		return models.MovieTimetable{}, ErrNotFound
	}

	return timetable, nil
}

func getScreenByID(screenID int) (models.Screen, error) {
	var screen models.Screen
	err := config.DB.QueryRow(context.Background(),
		"SELECT id, screen_no, normal, vip, type FROM screens WHERE id=$1",
		screenID,
	).Scan(&screen.ID, &screen.ScreenNo, &screen.Normal, &screen.Vip, &screen.Type)
	if err != nil {
		return models.Screen{}, ErrNotFound
	}

	return screen, nil
}

func getBookedSeats(timetableID, screenID int, dateTime time.Time) (map[string]struct{}, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT reservation FROM bookings WHERE timetable_id=$1 AND screen_id=$2 AND date_and_time=$3",
		timetableID, screenID, dateTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	booked := make(map[string]struct{})
	for rows.Next() {
		var seats []string
		if err := rows.Scan(&seats); err != nil {
			return nil, err
		}
		for _, seat := range seats {
			booked[seat] = struct{}{}
		}
	}

	return booked, nil
}

func normalizeSeats(seats []string) (map[string]struct{}, error) {
	unique := make(map[string]struct{}, len(seats))
	for _, seat := range seats {
		if seat == "" {
			return nil, errors.New("seat cannot be empty")
		}
		if _, exists := unique[seat]; exists {
			return nil, errors.New("duplicate seat found")
		}
		unique[seat] = struct{}{}
	}
	return unique, nil
}

func validateSeatsExist(seats map[string]struct{}, normal []string, vip []string) error {
	allowed := toStringSet(normal)
	for _, seat := range vip {
		allowed[seat] = struct{}{}
	}

	for seat := range seats {
		if _, exists := allowed[seat]; !exists {
			return errors.New("seat not found in screen")
		}
	}
	return nil
}

func containsInt(values []int, target int) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func toStringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func matchesTiming(dateTime time.Time, timings []string) bool {
	shortTime := dateTime.Format("15:04")
	fullTime := dateTime.Format("15:04:05")
	for _, timing := range timings {
		if timing == shortTime || timing == fullTime {
			return true
		}
	}
	return false
}

func isSameDate(left time.Time, right time.Time) bool {
	return left.Year() == right.Year() && left.Month() == right.Month() && left.Day() == right.Day()
}
