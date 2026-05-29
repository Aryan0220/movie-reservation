package services

import (
	"booking-system/config"
	"booking-system/models"
	"context"
)

type SeatStatusRequest struct {
	ShowTimeID int  `json:"showtime_id"`
	ScreenID   int	`json:"screen_id"`
	ShowTime   string `json:"show_time"`
}

func Add_Movie(movie models.Movie) error {
	_, err := config.DB.Exec(context.Background(),
		"INSERT INTO movies (title, description, poster_url, genre) VALUES ($1, $2, $3, $4)",
		movie.Title, movie.Description, movie.PosterURL, movie.Genre,
	)
	return err
}

func Update_Movie(movie models.Movie) error {
	_, err := config.DB.Exec(context.Background(),
		"UPDATE movies SET title=$1, description=$2, poster_url=$3, genre=$4 WHERE title=$5",
		movie.Title, movie.Description, movie.PosterURL, movie.Genre, movie.Title,
	)
	return err
}

func Delete_Movie(movie models.Movie) error {
	_, err := config.DB.Exec(context.Background(),
		"DELETE FROM movies WHERE title = $1",
		movie.Title,
	)
	return err
}

func Get_Movies() ([]string, error) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT title FROM MOVIES",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []string
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return nil, err
		}
		movies = append(movies, title)
	}

	return movies, nil
}

func ViewSeats(request SeatStatusRequest) map[string]string {
	status := config.DB.QueryRow(context.Background(),
		"SELECT seat_status FROM show_seats WHERE showtime_id=$1 AND screen_id=$2 AND show_time=$3",
		request.ShowTimeID, request.ScreenID, request.ShowTime)

	var result map[string]string
	err := status.Scan(&result)

	if err != nil {
		return nil
	}

	return result
}
