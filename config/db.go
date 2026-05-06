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
		return err
	}
	log.Println("connected to DB")
	return nil
}
