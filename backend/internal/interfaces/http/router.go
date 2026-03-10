package http

import (
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
