package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBConn interface {
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Begin(context.Context) (pgx.Tx, error)
}

var DB DBConn

func ConnectDB() error {
	var err error
	pool, err := pgxpool.New(context.Background(), GetEnv("DB_URL"))
	if err != nil {
		log.Printf("Failed to connect to DB: %v", err)
		return err
	}
	DB = pool
	log.Println("Connected to DB")
	return nil
}
