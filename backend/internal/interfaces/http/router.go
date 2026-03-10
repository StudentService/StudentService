package http

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs" // подключаем сгенерированные docs

	"backend/internal/application"
	"backend/internal/infrastructure"
	"backend/internal/interfaces/http/handlers"
	"backend/internal/interfaces/http/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Настройка CORS
	r.Use(cors.New(cors.Config{
		// Разрешенные origins (можно указать несколько)
		AllowOrigins: []string{
			"http://localhost:5173", // Vite dev server
			"http://127.0.0.1:5173",
			"http://localhost:3000", // Если вдруг другой порт
			"http://127.0.0.1:3000",
		},
		// Разрешенные методы
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		// Разрешенные заголовки
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		// Разрешаем отправку учетных данных (cookies, авторизационные заголовки)
		AllowCredentials: true,
		// Максимальное время кеширования preflight запросов
		MaxAge: 12 * time.Hour,
	}))

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Репозитории
	userRepo := &infrastructure.UserRepository{}

	// Сервисы
	userService := application.NewUserService(userRepo)
	authService := application.NewAuthService(userRepo)

	// Хендлеры
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Публичные маршруты (без аутентификации)
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	// Защищенные маршруты (требуют JWT)
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// Профиль текущего пользователя
		api.GET("/users/me", userHandler.GetMe)
		api.PATCH("/users/me", userHandler.UpdateMe)

		// Доступ к другим пользователям (с проверкой прав)
		api.GET("/users/:id", userHandler.GetUserByID)
	}

	return r
}
