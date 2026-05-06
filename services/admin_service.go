package services

import (
	"booking-system/config"
	"context"
	"booking-system/models"
)

func PromoteToAdmin(user models.User) error {
	if user.Role == "admin" {
		return nil
	}
	
	_, err := config.DB.Exec(context.Background(),
		"UPDATE users SET role='admin' WHERE email=$1",
		(user.Email),
	)

	return err
}

/*
func MovieRevenue(movie models.Movie) error {

}

func AllReservation() error {

}

func TotalRevenue() error {

}

func Capacity() error {

}
*/
