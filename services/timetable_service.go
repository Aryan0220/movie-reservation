package services

import (
	"booking-system/config"
	"booking-system/models"
	"context"
)

func AddShowTime(timetable models.MovieTimetable) error {
	_, err := config.DB.Exec(context.Background(),
		"INSERT INTO movie_timetables (movie_id, timings, screens, show_date, normal_price, vip_price) VALUES ($1, $2, $3, $4, $5, $6)",
		timetable.MovieID, timetable.Timings, timetable.Screens, timetable.ShowDate, timetable.NormalPrice, timetable.VipPrice,
	)
	return err
}

func UpdateShowTime(timetable models.MovieTimetable) error {
	_, err := config.DB.Exec(context.Background(),
		"UPDATE movie_timetables SET movie_id=$1, timings=$2, screens=$3, show_date=$4, normal_price=$5, vip_price=$6 WHERE id=$7",
		timetable.MovieID, timetable.Timings, timetable.Screens, timetable.ShowDate, timetable.NormalPrice, timetable.VipPrice, timetable.ID,
	)
	return err
}
