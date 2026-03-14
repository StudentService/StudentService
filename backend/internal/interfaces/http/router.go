package http

import (
	"backend/internal/infrastructure/repositories"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "backend/docs"

	"backend/internal/application"
	"backend/internal/interfaces/http/handlers"
	"backend/internal/interfaces/http/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:3000",
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Репозитории
	userRepo := &repositories.UserRepository{}
	semesterRepo := &repositories.SemesterRepository{}
	courseRepo := &repositories.CourseRepository{}
	groupRepo := &repositories.GroupRepository{}
	calendarRepo := &repositories.CalendarRepository{}
	challengeRepo := &repositories.ChallengeRepository{}
	gradeRepo := repositories.NewGradeRepository()
	questionnaireRepo := &repositories.QuestionnaireRepository{}
	activityRepo := repositories.NewActivityRepository()

	// Сервисы
	userService := application.NewUserService(userRepo)
	authService := application.NewAuthService(userRepo)
	calendarService := application.NewCalendarService(calendarRepo, userRepo, courseRepo, groupRepo)
	challengeService := application.NewChallengeService(challengeRepo, userRepo)
	gradeService := application.NewGradeService(gradeRepo, userRepo, courseRepo)
	questionnaireService := application.NewQuestionnaireService(questionnaireRepo, userRepo)
	activityService := application.NewActivityService(activityRepo, userRepo, courseRepo, groupRepo)
	dashboardService := application.NewDashboardService(
		userRepo,
		calendarRepo,
		challengeRepo,
		gradeRepo,
		questionnaireRepo,
		activityRepo,
		groupRepo,
		courseRepo,
		semesterRepo,
	)

	// Хендлеры
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	calendarHandler := handlers.NewCalendarHandler(calendarService)
	challengeHandler := handlers.NewChallengeHandler(challengeService)
	gradeHandler := handlers.NewGradeHandler(gradeService)
	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireService)
	activityHandler := handlers.NewActivityHandler(activityService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

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
		// Дашборд (главная страница)
		api.GET("/dashboard", dashboardHandler.GetStudentDashboard)

		// Профиль текущего пользователя
		api.GET("/users/me", userHandler.GetMe)
		api.PATCH("/users/me", userHandler.UpdateMe)

		// Доступ к другим пользователям (с проверкой прав)
		api.GET("/users/:id", userHandler.GetUserByID)

		// Календарь
		calendar := api.Group("/calendar")
		{
			calendar.GET("/events/my", calendarHandler.GetMyEvents)
			calendar.POST("/events", calendarHandler.CreateEvent)       // для преподавателей/админов
			calendar.PATCH("/events/:id", calendarHandler.UpdateEvent)  // для создателя/админа
			calendar.DELETE("/events/:id", calendarHandler.DeleteEvent) // для создателя/админа
		}

		// Личные вызовы
		api.GET("/challenges/my", challengeHandler.GetMyChallenges)
		api.GET("/challenges/:id", challengeHandler.GetChallengeByID)
		api.POST("/challenges", challengeHandler.CreateChallenge)
		api.PATCH("/challenges/:id", challengeHandler.UpdateChallenge)
		api.DELETE("/challenges/:id", challengeHandler.DeleteChallenge)

		// Оценки и успеваемость
		grades := api.Group("/grades")
		{
			// Для студента
			grades.GET("/my", gradeHandler.GetMyGrades)
			grades.GET("/my/summary", gradeHandler.GetMySummary)
			grades.GET("/my/period", gradeHandler.GetMyGradesByPeriod)
			grades.GET("/my/courses/:courseId", gradeHandler.GetMyGradesByCourse)

			// Для преподавателя (управление)
			grades.POST("/students/:studentId", gradeHandler.CreateGrade)
			grades.PATCH("/:id", gradeHandler.UpdateGrade)
			grades.DELETE("/:id", gradeHandler.DeleteGrade)
		}

		// Анкета-запрос
		questionnaire := api.Group("/questionnaire")
		{
			questionnaire.GET("/my", questionnaireHandler.GetMyQuestionnaire)
			questionnaire.GET("/template", questionnaireHandler.GetTemplate)
			questionnaire.POST("/submit", questionnaireHandler.SubmitQuestionnaire)
			questionnaire.POST("/draft", questionnaireHandler.SaveDraft)

			// Админские маршруты
			questionnaire.GET("", questionnaireHandler.ListByStatus)
			questionnaire.POST("/:id/review", questionnaireHandler.ReviewQuestionnaire)
		}

		// Активности и участия
		activities := api.Group("/activities")
		{
			// Для студентов
			activities.GET("/available", activityHandler.GetAvailableActivities)
			activities.GET("/my", activityHandler.GetMyParticipations)
			activities.POST("/:activityId/enroll", activityHandler.Enroll)
			activities.DELETE("/:activityId/enroll", activityHandler.CancelEnrollment)

			// Для преподавателей/админов - меняем :id на :activityId для консистентности
			activities.POST("", activityHandler.CreateActivity)
			activities.PATCH("/:activityId", activityHandler.UpdateActivity)                     // было :id, стало :activityId
			activities.DELETE("/:activityId", activityHandler.DeleteActivity)                    // было :id, стало :activityId
			activities.GET("/:activityId/participants", activityHandler.GetActivityParticipants) // было :id, стало :activityId
		}
	}

	return r
}
