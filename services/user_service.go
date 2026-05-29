package services

import (
	"context"
	"booking-system/config"
	"booking-system/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func CreateUser(user models.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	log.Printf("Creating user: %s with email: %s and role: %v", user.Name, user.Email, user.Role)
	_, err := config.DB.Exec(context.Background(),
		"INSERT INTO users (name, email, admin, password) VALUES ($1, $2, $3, $4)",
		user.Name, user.Email, user.Role, string(hash),
	)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := config.DB.QueryRow(context.Background(),
		"SELECT id, name, email, admin, password FROM users WHERE email=$1",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Password)

	return user, err
}

func GetMovieTimings(title string) ([]models.MovieTimetable, error) {
	var timings []models.MovieTimetable

	rows, err := config.DB.Query(context.Background(),
		"SELECT id, movie_id, schedule, show_date, normal_price, vip_price FROM showtimes WHERE movie_id=(SELECT id FROM movies WHERE title=$1)",
		title,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var timing models.MovieTimetable

		err := rows.Scan(&timing.ID, &timing.MovieID, &timing.Schedule, &timing.ShowDate, &timing.NormalPrice, &timing.VipPrice)

		if err != nil {
			return nil, err
		}

		timings = append(timings, timing)

	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}


	return timings, err
}
