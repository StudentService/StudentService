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

func SetupRouter(rbacService *application.RBACService) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5174",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger
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
	userService := application.NewUserService(userRepo, groupRepo, semesterRepo)
	groupService := application.NewGroupService(groupRepo)
	semesterService := application.NewSemesterService(semesterRepo)
	teacherService := application.NewTeacherService(userRepo, groupRepo, gradeRepo, activityRepo)
	authService := application.NewAuthService(userRepo)
	calendarService := application.NewCalendarService(calendarRepo, userRepo, courseRepo, groupRepo)
	challengeService := application.NewChallengeService(challengeRepo, userRepo)
	gradeService := application.NewGradeService(gradeRepo, userRepo, courseRepo)
	questionnaireService := application.NewQuestionnaireService(questionnaireRepo, userRepo)
	activityService := application.NewActivityService(activityRepo, userRepo, courseRepo, groupRepo)
	dashboardService := application.NewDashboardService(
		userRepo, calendarRepo, challengeRepo, gradeRepo,
		questionnaireRepo, activityRepo, groupRepo, courseRepo, semesterRepo,
	)

	// Хендлеры
	userHandler := handlers.NewUserHandler(userService)
	groupHandler := handlers.NewGroupHandler(groupService)
	semesterHandler := handlers.NewSemesterHandler(semesterService)
	teacherHandler := handlers.NewTeacherHandler(teacherService)
	authHandler := handlers.NewAuthHandler(authService)
	calendarHandler := handlers.NewCalendarHandler(calendarService)
	challengeHandler := handlers.NewChallengeHandler(challengeService)
	gradeHandler := handlers.NewGradeHandler(gradeService)
	questionnaireHandler := handlers.NewQuestionnaireHandler(questionnaireService)
	activityHandler := handlers.NewActivityHandler(activityService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	// Публичные маршруты
	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	// 👇 Публичные маршруты для групп (без токена)
	public := r.Group("/api/v1")
	{
		// Группы - только чтение (для регистрации)
		public.GET("/groups", groupHandler.GetAllGroups)
		public.GET("/groups/:id", groupHandler.GetGroupByID)

		// Семестры - только чтение (опционально)
		public.GET("/semesters", semesterHandler.GetAllSemesters)
		public.GET("/semesters/active", semesterHandler.GetActiveSemester)
	}

	// Защищённые маршруты
	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())
	{
		// Дашборд
		api.GET("/dashboard",
			middleware.RequirePermission(rbacService, "dashboard", "read"),
			dashboardHandler.GetStudentDashboard,
		)

		// Профиль
		api.GET("/users/me",
			middleware.RequirePermission(rbacService, "profile", "read"),
			userHandler.GetMe,
		)
		api.PATCH("/users/me",
			middleware.RequirePermission(rbacService, "profile", "write"),
			userHandler.UpdateMe,
		)
		api.GET("/users/:id",
			middleware.RequirePermission(rbacService, "profile", "read"),
			userHandler.GetUserByID,
		)

		// CRUD операции для групп (требуют аутентификации)
		api.POST("/groups", groupHandler.CreateGroup)
		api.PATCH("/groups/:id", groupHandler.UpdateGroup)
		api.DELETE("/groups/:id", groupHandler.DeleteGroup)

		// CRUD операции для семестров
		api.POST("/semesters", semesterHandler.CreateSemester)
		api.PATCH("/semesters/:id", semesterHandler.UpdateSemester)
		api.DELETE("/semesters/:id", semesterHandler.DeleteSemester)

		// Календарь
		calendar := api.Group("/calendar")
		{
			calendar.GET("/events/my",
				middleware.RequirePermission(rbacService, "calendar", "read"),
				calendarHandler.GetMyEvents,
			)
			calendar.POST("/events",
				middleware.RequirePermission(rbacService, "calendar", "write"),
				calendarHandler.CreateEvent,
			)
			calendar.PATCH("/events/:id",
				middleware.RequirePermission(rbacService, "calendar", "write"),
				calendarHandler.UpdateEvent,
			)
			calendar.DELETE("/events/:id",
				middleware.RequirePermission(rbacService, "calendar", "write"),
				calendarHandler.DeleteEvent,
			)
		}

		// Личные вызовы
		challenges := api.Group("/challenges")
		{
			challenges.GET("/my",
				middleware.RequirePermission(rbacService, "challenge", "read"),
				challengeHandler.GetMyChallenges,
			)
			challenges.GET("/:id",
				middleware.RequirePermission(rbacService, "challenge", "read"),
				challengeHandler.GetChallengeByID,
			)
			challenges.POST("",
				middleware.RequirePermission(rbacService, "challenge", "write"),
				challengeHandler.CreateChallenge,
			)
			challenges.PATCH("/:id",
				middleware.RequirePermission(rbacService, "challenge", "write"),
				challengeHandler.UpdateChallenge,
			)
			challenges.DELETE("/:id",
				middleware.RequirePermission(rbacService, "challenge", "delete"),
				challengeHandler.DeleteChallenge,
			)
		}

		// Оценки
		grades := api.Group("/grades")
		{
			grades.GET("/my",
				middleware.RequirePermission(rbacService, "grade", "read"),
				gradeHandler.GetMyGrades,
			)
			grades.GET("/my/summary",
				middleware.RequirePermission(rbacService, "grade", "read"),
				gradeHandler.GetMySummary,
			)
			grades.GET("/my/period",
				middleware.RequirePermission(rbacService, "grade", "read"),
				gradeHandler.GetMyGradesByPeriod,
			)
			grades.GET("/my/courses/:courseId",
				middleware.RequirePermission(rbacService, "grade", "read"),
				gradeHandler.GetMyGradesByCourse,
			)
			grades.POST("/students/:studentId",
				middleware.RequirePermission(rbacService, "grade", "write"),
				gradeHandler.CreateGrade,
			)
			grades.PATCH("/:id",
				middleware.RequirePermission(rbacService, "grade", "write"),
				gradeHandler.UpdateGrade,
			)
			grades.DELETE("/:id",
				middleware.RequirePermission(rbacService, "grade", "delete"),
				gradeHandler.DeleteGrade,
			)
		}

		// Анкета
		questionnaire := api.Group("/questionnaire")
		{
			questionnaire.GET("/my",
				middleware.RequirePermission(rbacService, "questionnaire", "read"),
				questionnaireHandler.GetMyQuestionnaire,
			)
			questionnaire.GET("/template",
				middleware.RequirePermission(rbacService, "questionnaire", "read"),
				questionnaireHandler.GetTemplate,
			)
			questionnaire.POST("/submit",
				middleware.RequirePermission(rbacService, "questionnaire", "write"),
				questionnaireHandler.SubmitQuestionnaire,
			)
			questionnaire.POST("/draft",
				middleware.RequirePermission(rbacService, "questionnaire", "write"),
				questionnaireHandler.SaveDraft,
			)
		}

		// Активности
		activities := api.Group("/activities")
		{
			activities.GET("/available",
				middleware.RequirePermission(rbacService, "activity", "read"),
				activityHandler.GetAvailableActivities,
			)
			activities.GET("/my",
				middleware.RequirePermission(rbacService, "activity", "read"),
				activityHandler.GetMyParticipations,
			)
			activities.POST("/:activityId/enroll",
				middleware.RequirePermission(rbacService, "activity", "enroll"),
				activityHandler.Enroll,
			)
			activities.DELETE("/:activityId/enroll",
				middleware.RequirePermission(rbacService, "activity", "enroll"),
				activityHandler.CancelEnrollment,
			)
			activities.POST("",
				middleware.RequirePermission(rbacService, "activity", "write"),
				activityHandler.CreateActivity,
			)
			activities.PATCH("/:activityId",
				middleware.RequirePermission(rbacService, "activity", "write"),
				activityHandler.UpdateActivity,
			)
			activities.DELETE("/:activityId",
				middleware.RequirePermission(rbacService, "activity", "delete"),
				activityHandler.DeleteActivity,
			)
			activities.GET("/:activityId/participants",
				middleware.RequirePermission(rbacService, "activity", "read"),
				activityHandler.GetActivityParticipants,
			)
		}

		// Преподавательские маршруты
		teacher := api.Group("/teacher")
		teacher.Use(middleware.RequirePermission(rbacService, "teacher", "access"))
		{
			teacher.GET("/dashboard", teacherHandler.GetDashboard)
			teacher.GET("/groups", teacherHandler.GetTeacherGroups)
			teacher.GET("/groups/:id/students", teacherHandler.GetGroupStudents)
			teacher.GET("/groups/:id/grades", teacherHandler.GetGroupGrades)
			teacher.GET("/students/:id", teacherHandler.GetStudentProfile)
			teacher.GET("/students/:id/grades", teacherHandler.GetStudentGrades)
			teacher.GET("/activities", teacherHandler.GetTeacherActivities)
			teacher.POST("/grades/import", teacherHandler.ImportGrades)
			teacher.POST("/activities/:id/attendance", teacherHandler.MarkAttendance)
		}
	}

	return r
}
