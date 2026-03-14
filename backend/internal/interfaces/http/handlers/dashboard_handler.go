package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
)

type DashboardHandler struct {
	dashboardService *application.DashboardService
}

func NewDashboardHandler(dashboardService *application.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetStudentDashboard godoc
// @Summary      Получение дашборда студента
// @Description  Возвращает сводную информацию для главной страницы студента
// @Tags         dashboard
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  dashboard.DashboardResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /dashboard [get]
func (h *DashboardHandler) GetStudentDashboard(c *gin.Context) {
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

	// Проверяем, что пользователь - студент (опционально)
	role, exists := c.Get("user_role")
	if exists && role != "student" {
		// Не студент, но может иметь доступ? Пока разрешим всем
		// Можно добавить логику
	}

	dashboard, err := h.dashboardService.GetStudentDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}
