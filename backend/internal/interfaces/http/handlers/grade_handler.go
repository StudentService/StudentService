package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/grade"
)

type GradeHandler struct {
	gradeService *application.GradeService
}

func NewGradeHandler(gradeService *application.GradeService) *GradeHandler {
	return &GradeHandler{gradeService: gradeService}
}

// GetMyGrades godoc
// @Summary      Получение всех оценок текущего студента
// @Description  Возвращает список всех оценок текущего пользователя
// @Tags         grades
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   grade.GradeResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /grades/my [get]
func (h *GradeHandler) GetMyGrades(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	grades, err := h.gradeService.GetMyGrades(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetMyGradesByCourse godoc
// @Summary      Получение оценок по курсу
// @Description  Возвращает оценки текущего студента по конкретному курсу
// @Tags         grades
// @Security     BearerAuth
// @Produce      json
// @Param        courseId path string true "Course ID"
// @Success      200  {array}   grade.GradeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /grades/my/courses/{courseId} [get]
func (h *GradeHandler) GetMyGradesByCourse(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	courseID, err := uuid.Parse(c.Param("courseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course id"})
		return
	}

	grades, err := h.gradeService.GetMyGradesByCourse(c.Request.Context(), userID, courseID)
	if err != nil {
		if err.Error() == "course not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetMyGradesByPeriod godoc
// @Summary      Получение оценок за период
// @Description  Возвращает оценки текущего студента за указанный период
// @Tags         grades
// @Security     BearerAuth
// @Produce      json
// @Param        from query string true "Начало периода (RFC3339)"
// @Param        to query string true "Конец периода (RFC3339)"
// @Success      200  {array}   grade.GradeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /grades/my/period [get]
func (h *GradeHandler) GetMyGradesByPeriod(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to parameters are required"})
		return
	}

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format"})
		return
	}

	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format"})
		return
	}

	grades, err := h.gradeService.GetMyGradesByPeriod(c.Request.Context(), userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, grades)
}

// GetMySummary godoc
// @Summary      Получение сводки успеваемости
// @Description  Возвращает общую сводку по успеваемости студента
// @Tags         grades
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  grade.StudentSummary
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /grades/my/summary [get]
func (h *GradeHandler) GetMySummary(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	summary, err := h.gradeService.GetMySummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// CreateGrade godoc
// @Summary      Создание новой оценки (для преподавателей)
// @Description  Создаёт новую оценку для студента
// @Tags         grades
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        studentId path string true "Student ID"
// @Param        request body grade.CreateGradeRequest true "Данные оценки"
// @Success      201  {object}  grade.GradeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /grades/students/{studentId} [post]
func (h *GradeHandler) CreateGrade(c *gin.Context) {
	creatorIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	creatorID, err := uuid.Parse(creatorIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid student id"})
		return
	}

	var req grade.CreateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.gradeService.CreateGrade(c.Request.Context(), creatorID, studentID, &req)
	if err != nil {
		switch err.Error() {
		case "student not found", "course not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions to create grades":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, created)
}

// UpdateGrade godoc
// @Summary      Обновление оценки
// @Description  Обновляет существующую оценку
// @Tags         grades
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Grade ID"
// @Param        request body grade.UpdateGradeRequest true "Данные для обновления"
// @Success      200  {object}  grade.GradeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /grades/{id} [patch]
func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	gradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grade id"})
		return
	}

	var req grade.UpdateGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.gradeService.UpdateGrade(c.Request.Context(), userID, gradeID, &req)
	if err != nil {
		switch err.Error() {
		case "grade not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteGrade godoc
// @Summary      Удаление оценки
// @Description  Удаляет оценку
// @Tags         grades
// @Security     BearerAuth
// @Param        id path string true "Grade ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /grades/{id} [delete]
func (h *GradeHandler) DeleteGrade(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	gradeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grade id"})
		return
	}

	err = h.gradeService.DeleteGrade(c.Request.Context(), userID, gradeID)
	if err != nil {
		switch err.Error() {
		case "grade not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
