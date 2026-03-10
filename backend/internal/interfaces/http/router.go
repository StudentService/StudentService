package http

import (
	"backend/internal/application"
	"backend/internal/infrastructure"
	"backend/internal/interfaces/http/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	userRepo := &infrastructure.UserRepository{}
	userService := application.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	r.GET("/api/v1/profile", userHandler.GetProfile)

	return r
}
