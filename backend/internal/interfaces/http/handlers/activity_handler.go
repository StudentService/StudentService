package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/activity"
)

type ActivityHandler struct {
	activityService *application.ActivityService
}

func NewActivityHandler(activityService *application.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

// GetAvailableActivities godoc
// @Summary      Получение доступных активностей
// @Description  Возвращает список активностей, доступных для текущего пользователя
// @Tags         activities
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   activity.ActivityResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/available [get]
func (h *ActivityHandler) GetAvailableActivities(c *gin.Context) {
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

	activities, err := h.activityService.GetAvailableActivities(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// GetMyParticipations godoc
// @Summary      Получение моих участий
// @Description  Возвращает список активностей, в которых участвует текущий пользователь
// @Tags         activities
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   activity.ParticipationResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/my [get]
func (h *ActivityHandler) GetMyParticipations(c *gin.Context) {
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

	participations, err := h.activityService.GetMyParticipations(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participations)
}

// Enroll godoc
// @Summary      Запись на активность
// @Description  Записывает текущего пользователя на указанную активность
// @Tags         activities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        activityId path string true "Activity ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/{activityId}/enroll [post]
func (h *ActivityHandler) Enroll(c *gin.Context) {
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

	activityID, err := uuid.Parse(c.Param("activityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	err = h.activityService.Enroll(c.Request.Context(), userID, activityID)
	if err != nil {
		switch err.Error() {
		case "activity not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "activity is not active":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "already enrolled in this activity":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case "no available slots":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully enrolled"})
}

// CancelEnrollment godoc
// @Summary      Отмена записи на активность
// @Description  Отменяет запись текущего пользователя на активность
// @Tags         activities
// @Security     BearerAuth
// @Param        activityId path string true "Activity ID"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/{activityId}/enroll [delete]
func (h *ActivityHandler) CancelEnrollment(c *gin.Context) {
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

	activityID, err := uuid.Parse(c.Param("activityId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	err = h.activityService.CancelEnrollment(c.Request.Context(), userID, activityID)
	if err != nil {
		switch err.Error() {
		case "not enrolled in this activity":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "cannot cancel enrollment after activity has started":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "enrollment cancelled"})
}

// CreateActivity godoc
// @Summary      Создание новой активности (для преподавателей/админов)
// @Description  Создаёт новую активность в каталоге
// @Tags         activities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body activity.CreateActivityRequest true "Данные активности"
// @Success      201  {object}  activity.Activity
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities [post]
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	creatorID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req activity.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.activityService.CreateActivity(c.Request.Context(), creatorID, &req)
	if err != nil {
		if err.Error() == "insufficient permissions to create activities" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, a)
}

// UpdateActivity godoc
// @Summary      Обновление активности
// @Description  Обновляет существующую активность
// @Tags         activities
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Activity ID"
// @Param        request body activity.UpdateActivityRequest true "Данные для обновления"
// @Success      200  {object}  activity.Activity
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/{id} [patch]
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
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

	activityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var req activity.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	a, err := h.activityService.UpdateActivity(c.Request.Context(), userID, activityID, &req)
	if err != nil {
		switch err.Error() {
		case "activity not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, a)
}

// DeleteActivity godoc
// @Summary      Удаление активности
// @Description  Удаляет активность (только если нет участников)
// @Tags         activities
// @Security     BearerAuth
// @Param        id path string true "Activity ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/{id} [delete]
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
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

	activityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	err = h.activityService.DeleteActivity(c.Request.Context(), userID, activityID)
	if err != nil {
		switch err.Error() {
		case "activity not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case "cannot delete activity with enrolled participants":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// GetActivityParticipants godoc
// @Summary      Получение участников активности (для преподавателей)
// @Description  Возвращает список студентов, записанных на активность
// @Tags         activities
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Activity ID"
// @Success      200  {array}   activity.ParticipationResponse
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /activities/{id}/participants [get]
func (h *ActivityHandler) GetActivityParticipants(c *gin.Context) {
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

	activityID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	participants, err := h.activityService.GetActivityParticipants(c.Request.Context(), userID, activityID)
	if err != nil {
		switch err.Error() {
		case "activity not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, participants)
}
