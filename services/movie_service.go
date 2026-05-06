package services

import (
	"context"
	"booking-system/models"
	"booking-system/config"
)

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
		movie.Title, movie.Description, movie.PosterURL, movie.Genre, movie.ID,
	)
	return err
}

func Delete_Movie(movie models.Movie) error {
	_, err := config.DB.Exec(context.Background(), 
		"DELETE FROM movies WHERE title = $1",
		movie.ID,
	)
	return err
}
/*
func AssignShowtime() error {

}
*/
