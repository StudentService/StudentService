package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs"
)

// @title Платформа «Молодёжь» — Студенческая часть (прототип)
// @version 0.1.0
// @description API для студента. Доступны только свои данные. Все эндпоинты защищены Bearer JWT.
// @termsOfService http://example.com/terms/

// @contact.name Backend Team
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name Internal Use Only
// @license.url -

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Определяем маршруты
	api := r.Group("/api/v1")
	{
		api.POST("/auth/login", loginHandler)
		api.GET("/profile", profileHandler)
		api.GET("/challenges/my", myChallengesHandler)
	}

	log.Println("Swagger: http://localhost:8080/swagger/index.html")
	r.Run(":8080")
}

// loginHandler обрабатывает вход студента
// @Summary Вход студента
// @Description Авторизация и получение JWT-токена
// @Tags auth
// @Accept json
// @Produce json
// @Param body body object{login=string,password=string} true "Данные для входа"
// @Success 200 {object} object{token=string,user=object}
// @Failure 401 {object} object{message=string}
// @Router /auth/login [post]
func loginHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "mock login"})
}

// profileHandler возвращает профиль текущего студента
// @Summary Получить профиль текущего студента
// @Description Возвращает данные профиля
// @Tags profile
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{id=string,firstName=string,lastName=string,group=string,semester=string,segment=string}
// @Failure 401 {object} object{message=string}
// @Router /profile [get]
func profileHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "мой профиль"})
}

// myChallengesHandler возвращает список личных вызовов
// @Summary Список личных вызовов
// @Tags challenges
// @Security BearerAuth
// @Produce json
// @Success 200 {array} object{id=string,title=string,status=string,progress=number}
// @Router /challenges/my [get]
func myChallengesHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "список вызовов"})
}
