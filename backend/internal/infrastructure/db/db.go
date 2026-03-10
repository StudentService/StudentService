package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	// Берем строку подключения из переменной окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback для разработки (не рекомендуется, но для удобства)
		dsn = "postgres://user:pass@localhost:5431/db?sslmode=disable"
		log.Println("Using default database connection string")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Проверяем подключение
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	Pool = pool
	log.Println("Database connected successfully")
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}
