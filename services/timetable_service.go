package services

import (
	"booking-system/config"
	"booking-system/models"
	"context"
	"encoding/json"
	"strconv"
	"time"
)

func AddShowTime(timetable models.MovieTimetable) error {
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

	var showtimeID int
	err = transaction.QueryRow(ctx,
		"INSERT INTO showtimes (movie_id, schedule, show_date, normal_price, vip_price) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		timetable.MovieID, timetable.Schedule, timetable.ShowDate, timetable.NormalPrice, timetable.VipPrice,
	).Scan(&showtimeID)
	if err != nil {
		config.PrintLog("Failed to insert showtime for movie id "+strconv.Itoa(timetable.MovieID)+": "+err.Error(), "ERROR")
		return err
	}

	for _, schedule := range timetable.Schedule {
		screen, screenErr := getScreenByID(schedule.ScreenID)
		if screenErr != nil {
			err = screenErr
			return err
		}

		seatStatus := buildSeatStatus(screen.Normal, screen.Vip)
		seatStatusJSON, marshalErr := json.Marshal(seatStatus)
		if marshalErr != nil {
			err = marshalErr
			return err
		}

		for _, timing := range schedule.Timings {
			showTime := normalizeShowTime(timing)
			_, execErr := transaction.Exec(ctx,
				"INSERT INTO show_seats (showtime_id, screen_id, show_time, seat_status) VALUES ($1, $2, $3, $4::jsonb)",
				showtimeID, schedule.ScreenID, showTime, string(seatStatusJSON),
			)
			if execErr != nil {
				err = execErr
				return err
			}
		}
	}

	if commitErr := transaction.Commit(ctx); commitErr != nil {
		return commitErr
	}

	return nil
}

func UpdateShowTime(timetable models.MovieTimetable) error {
	_, err := config.DB.Exec(context.Background(),
		"UPDATE showtimes SET movie_id=$1, schedule=$2, show_date=$3, normal_price=$4, vip_price=$5 WHERE id=$6",
		timetable.MovieID, timetable.Schedule, timetable.ShowDate, timetable.NormalPrice, timetable.VipPrice, timetable.ID,
	)
	return err
}

func buildSeatStatus(normal []string, vip []string) map[string]string {
	status := make(map[string]string, len(normal)+len(vip))
	for _, seat := range normal {
		status[seat] = "available"
	}
	for _, seat := range vip {
		status[seat] = "available"
	}
	return status
}

func normalizeShowTime(value string) string {
	if parsed, err := time.Parse("15:04:05", value); err == nil {
		return parsed.Format("15:04")
	}
	if parsed, err := time.Parse("15:04", value); err == nil {
		return parsed.Format("15:04")
	}
	return value
}
