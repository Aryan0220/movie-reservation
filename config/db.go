package config

import(
	"context"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
) 

var DB *pgxpool.Pool

func ConnectDB() error {
	var err error
	DB, err = pgxpool.New(context.Background(), GetEnv("DB_URL"))
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		return err
	}
	log.Println("Connected to DB")
	return nil
}
