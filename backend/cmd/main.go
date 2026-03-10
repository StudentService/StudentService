package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"backend/internal/infrastructure/db"
	"backend/internal/interfaces/http"
)

// @title           Student Platform API
// @version         1.0
// @description     API для студенческой платформы Центра поддержки и развития молодёжи
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@studentplatform.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализация БД
	db.Init()
	defer db.Close()

	if err := db.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Database connected successfully")

	// Настройка и запуск сервера
	r := http.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
