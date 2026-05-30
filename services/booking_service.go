package services

import (
	"booking-system/config"
	"booking-system/models"
	"context"
	"encoding/json"
	"errors"
	"time"
	"log"
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


	if !isSameDate(timetable.ShowDate, dateTime) || !matchesTiming(dateTime, getTimingsByScreenID(screenID, timetable.Schedule)) {
		return ErrInvalidShowtime
	}

	if getTimingsByScreenID(screenID, timetable.Schedule) == nil {
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

	if !checkSeatAvailability(timetableID, screenID, dateTime, requestedSeats) {
		return ErrSeatUnavailable
	}

	_, err = config.DB.Exec(context.Background(),
		"INSERT INTO bookings (user_id, showtime_id, screen_id, seats, booking_time) VALUES ($1, $2, $3, $4, $5)",
		userID, timetableID, screenID, seats, dateTime,
	)
	if err != nil {
		return err
	}

	seatPatch, patchErr := buildSeatStatusPatch(requestedSeats, "booked")
	if patchErr != nil {
		return patchErr
	}

	commandTag, updateErr := config.DB.Exec(context.Background(),
		"UPDATE show_seats SET seat_status = seat_status || $4::jsonb WHERE showtime_id=$1 AND screen_id=$2 AND show_time=$3",
		timetableID, screenID, dateTime.Format("15:04"), seatPatch,
	)
	if updateErr != nil {
		return updateErr
	}
	if commandTag.RowsAffected() == 0 {
		return ErrInvalidShowtime
	}

	return nil
}

func CancelReservation(userID, bookingID int, now time.Time) error {
	ctx := context.Background()
	transaction, err := config.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = transaction.Rollback(ctx)
		}
	}()

	var ownerID int
	var timetableID int
	var screenID int
	var seats []string
	var dateTime time.Time

	err = transaction.QueryRow(ctx,
		"SELECT user_id, showtime_id, screen_id, seats, booking_time FROM bookings WHERE id=$1",
		bookingID,
	).Scan(&ownerID, &timetableID, &screenID, &seats, &dateTime)
	if err != nil {
		return ErrNotFound
	}

	if ownerID != userID {
		return ErrNotOwner
	}

	if !dateTime.After(now) {
		return ErrPastReservation
	}

	seatSet := toStringSet(seats)
	seatPatch, patchErr := buildSeatStatusPatch(seatSet, "available")
	if patchErr != nil {
		err = patchErr
		return err
	}

	commandTag, updateErr := transaction.Exec(ctx,
		"UPDATE show_seats SET seat_status = seat_status || $4::jsonb WHERE showtime_id=$1 AND screen_id=$2 AND show_time=$3",
		timetableID, screenID, dateTime.Format("15:04"), seatPatch,
	)
	if updateErr != nil {
		err = updateErr
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrInvalidShowtime
	}

	_, err = transaction.Exec(ctx,
		"DELETE FROM bookings WHERE id=$1",
		bookingID,
	)
	if err != nil {
		return err
	}

	if commitErr := transaction.Commit(ctx); commitErr != nil {
		return commitErr
	}

	return nil
}

func GetCapacity(timetableID, screenID int, dateTime time.Time) (int, int, error) {
	timetable, err := getTimetableByID(timetableID)
	if err != nil {
		return 0, 0, err
	}

	if !isSameDate(timetable.ShowDate, dateTime) || !matchesTiming(dateTime, getTimingsByScreenID(screenID, timetable.Schedule))  {
		return 0, 0, ErrInvalidShowtime
	}

	if getTimingsByScreenID(screenID, timetable.Schedule) == nil  {
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
		"SELECT id, user_id, showtime_id, screen_id, seats, booking_time FROM bookings ORDER BY booking_time DESC",
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
		"SELECT b.screen_id, b.seats, mt.normal_price, mt.vip_price FROM bookings b JOIN showtimes mt ON mt.id = b.showtime_id WHERE mt.movie_id=$1",
		movieID)
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
		"SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes WHERE id=$1",
		timetableID,
	).Scan(&timetable.ID, &timetable.MovieID, &timetable.Schedule, &timetable.ShowDate, &timetable.NormalPrice, &timetable.VipPrice)
	if err != nil {
		return models.MovieTimetable{}, ErrNotFound
	}

	return timetable, nil
}

func getScreenByID(screenID int) (models.Screen, error) {
	var screen models.Screen
	err := config.DB.QueryRow(context.Background(),
		"SELECT id, auditorium_number, normal_seats, vip_seats, type FROM screens WHERE id=$1",
		screenID,
	).Scan(&screen.ID, &screen.ScreenNo, &screen.Normal, &screen.Vip, &screen.Type)
	if err != nil {
		return models.Screen{}, ErrNotFound
	}

	return screen, nil
}

func getBookedSeats(timetableID, screenID int, dateTime time.Time) (map[string]struct{}, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT seats FROM bookings WHERE showtime_id=$1 AND screen_id=$2 AND booking_time=$3",
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

func toStringSet(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func getTimingsByScreenID(screenID int, data []models.Schedule) []string {
	for _, item := range data {
		if item.ScreenID == screenID {
			return item.Timings
		}
	}
	return nil
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

func checkSeatAvailability(timetableID, screenID int, dateTime time.Time, requestedSeats map[string]struct{}) bool {
	status := config.DB.QueryRow(context.Background(),
		"SELECT seat_status FROM show_seats WHERE showtime_id=$1 AND screen_id=$2 AND show_time=$3",
		timetableID, screenID, dateTime.Format("15:04"))

	var seatStatus map[string]string
	err := status.Scan(&seatStatus)
	if err != nil {
		log.Printf("Failed to check seat availability: %v", err)
		return false
	}
	log.Println(seatStatus, requestedSeats)
	for seat := range requestedSeats {
		if status, exists := seatStatus[seat]; !exists || status != "available" {
			return false
		}
	}
	return true
}

func buildSeatStatusPatch(seats map[string]struct{}, status string) (string, error) {
	patch := make(map[string]string, len(seats))
	for seat := range seats {
		patch[seat] = status
	}
	encoded, err := json.Marshal(patch)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}