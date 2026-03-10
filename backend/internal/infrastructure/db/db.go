package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	dsn := "postgres://user:pass@localhost:5431/db?sslmode=disable" // из env в проде
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	Pool = pool
}

func Close() {
	Pool.Close()
}
