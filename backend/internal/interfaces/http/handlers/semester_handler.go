package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/semester"
)

type SemesterHandler struct {
	semesterService *application.SemesterService
}

func NewSemesterHandler(semesterService *application.SemesterService) *SemesterHandler {
	return &SemesterHandler{semesterService: semesterService}
}

// GetAllSemesters godoc
// @Summary      Получение всех семестров
// @Tags         semesters
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  semester.Semester
// @Router       /semesters [get]
func (h *SemesterHandler) GetAllSemesters(c *gin.Context) {
	semesters, err := h.semesterService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, semesters)
}

// GetActiveSemester godoc
// @Summary      Получение активного семестра
// @Tags         semesters
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  semester.Semester
// @Router       /semesters/active [get]
func (h *SemesterHandler) GetActiveSemester(c *gin.Context) {
	s, err := h.semesterService.GetActive(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if s == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no active semester"})
		return
	}
	c.JSON(http.StatusOK, s)
}

// CreateSemester godoc
// @Summary      Создание семестра
// @Tags         semesters
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body semester.CreateSemesterRequest true "Данные семестра"
// @Success      201  {object}  semester.Semester
// @Router       /semesters [post]
func (h *SemesterHandler) CreateSemester(c *gin.Context) {
	var req semester.CreateSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := h.semesterService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

// UpdateSemester godoc
// @Summary      Обновление семестра
// @Tags         semesters
// @Security     BearerAuth
// @Param        id path string true "Semester ID"
// @Param        request body semester.UpdateSemesterRequest true "Данные для обновления"
// @Success      200  {object}  semester.Semester
// @Router       /semesters/{id} [patch]
func (h *SemesterHandler) UpdateSemester(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid semester id"})
		return
	}

	var req semester.UpdateSemesterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := h.semesterService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// DeleteSemester godoc
// @Summary      Удаление семестра
// @Tags         semesters
// @Security     BearerAuth
// @Param        id path string true "Semester ID"
// @Success      204  "No Content"
// @Router       /semesters/{id} [delete]
func (h *SemesterHandler) DeleteSemester(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid semester id"})
		return
	}

	if err := h.semesterService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
