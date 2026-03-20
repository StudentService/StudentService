package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/teacher"
)

type TeacherHandler struct {
	teacherService *application.TeacherService
}

func NewTeacherHandler(teacherService *application.TeacherService) *TeacherHandler {
	return &TeacherHandler{teacherService: teacherService}
}

// GetDashboard godoc
// @Summary      Дашборд преподавателя
// @Description  Возвращает сводку по группам и активности
// @Tags         teacher
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  teacher.DashboardResponse
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Router       /teacher/dashboard [get]
func (h *TeacherHandler) GetDashboard(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	dashboard, err := h.teacherService.GetDashboard(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

// GetTeacherGroups godoc
// @Summary      Группы преподавателя
// @Description  Возвращает список групп преподавателя
// @Tags         teacher
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   teacher.GroupSummary
// @Router       /teacher/groups [get]
func (h *TeacherHandler) GetTeacherGroups(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	groups, err := h.teacherService.GetTeacherGroups(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// GetGroupStudents godoc
// @Summary      Студенты группы
// @Description  Возвращает список студентов в группе
// @Tags         teacher
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Success      200  {array}   user.UserResponse
// @Router       /teacher/groups/{id}/students [get]
func (h *TeacherHandler) GetGroupStudents(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	students, err := h.teacherService.GetGroupStudents(c.Request.Context(), teacherID, groupID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, students)
}

// GetStudentProfile godoc
// @Summary      Профиль студента
// @Description  Возвращает профиль студента для преподавателя
// @Tags         teacher
// @Security     BearerAuth
// @Param        id path string true "Student ID"
// @Success      200  {object}  teacher.StudentProfile
// @Router       /teacher/students/{id} [get]
func (h *TeacherHandler) GetStudentProfile(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	profile, err := h.teacherService.GetStudentProfile(c.Request.Context(), teacherID, studentID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetStudentGrades godoc
// @Summary      Оценки студента
// @Description  Возвращает оценки студента
// @Tags         teacher
// @Security     BearerAuth
// @Param        id path string true "Student ID"
// @Success      200  {array}   teacher.GradeWithStudent
// @Router       /teacher/students/{id}/grades [get]
func (h *TeacherHandler) GetStudentGrades(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	studentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	grades, err := h.teacherService.GetStudentGrades(c.Request.Context(), teacherID, studentID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetGroupGrades godoc
// @Summary      Оценки группы
// @Description  Возвращает все оценки студентов группы
// @Tags         teacher
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Success      200  {array}   teacher.GradeWithStudent
// @Router       /teacher/groups/{id}/grades [get]
func (h *TeacherHandler) GetGroupGrades(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	grades, err := h.teacherService.GetGroupGrades(c.Request.Context(), teacherID, groupID)
	if err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetTeacherActivities godoc
// @Summary      Активности преподавателя
// @Description  Возвращает список активностей, созданных преподавателем
// @Tags         teacher
// @Security     BearerAuth
// @Success      200  {array}   activity.ActivityResponse
// @Router       /teacher/activities [get]
func (h *TeacherHandler) GetTeacherActivities(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	activities, err := h.teacherService.GetTeacherActivities(c.Request.Context(), teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// ImportGrades godoc
// @Summary      Импорт оценок
// @Description  Импортирует оценки из CSV/JSON
// @Tags         teacher
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body teacher.ImportGradesRequest true "Данные для импорта"
// @Success      200  {object}  map[string]string
// @Router       /teacher/grades/import [post]
func (h *TeacherHandler) ImportGrades(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req teacher.ImportGradesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.teacherService.ImportGrades(c.Request.Context(), teacherID, &req); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grades imported successfully"})
}

// MarkAttendance godoc
// @Summary      Отметка посещаемости
// @Description  Отмечает посещаемость студентов на активности
// @Tags         teacher
// @Security     BearerAuth
// @Param        id path string true "Activity ID"
// @Param        request body teacher.AttendanceRequest true "Данные посещаемости"
// @Success      200  {object}  map[string]string
// @Router       /teacher/activities/{id}/attendance [post]
func (h *TeacherHandler) MarkAttendance(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	teacherID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	activityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var req teacher.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.teacherService.MarkAttendance(c.Request.Context(), teacherID, activityID, &req); err != nil {
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "attendance marked successfully"})
}
