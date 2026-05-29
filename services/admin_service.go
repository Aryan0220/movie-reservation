package services

import (
	"booking-system/config"
	"context"
	"booking-system/models"
)

func PromoteToAdmin(user models.User) error {
	if user.Role == true{
		return nil
	}
	
	_, err := config.DB.Exec(context.Background(),
		"UPDATE users SET admin=true WHERE email=$1",
		(user.Email),
	)

	return err
}

 