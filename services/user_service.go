package services

import (
	"context"
	"booking-system/config"
	"booking-system/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user models.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	_, err := config.DB.Exec(context.Background(),
		"INSERT INTO users (name, email, role, password) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Role, string(hash),
	)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := config.DB.QueryRow(context.Background(),
		"SELECT id, name, email, role, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)

	return user, err
}

func GetMovieTimings(title string) (models.MovieTimetable, error) {
	var timings models.MovieTimetable

	err := config.DB.QueryRow(context.Background(),
		"SELECT id, movie_id, timings, screens, show_date, normal_price, vip_price FROM movie_timetables WHERE movie_id=(SELECT id FROM movies WHERE title=$1)",
		title,
	).Scan(&timings.ID, &timings.MovieID, &timings.Timings, &timings.Screens, &timings.ShowDate, &timings.NormalPrice, &timings.VipPrice)

	return timings, err
}
