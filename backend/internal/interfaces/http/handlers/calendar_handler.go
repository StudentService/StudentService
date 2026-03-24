package handlers

import (
	"backend/internal/domain/calendar"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
)

type CalendarHandler struct {
	calendarService *application.CalendarService
}

func NewCalendarHandler(calendarService *application.CalendarService) *CalendarHandler {
	return &CalendarHandler{calendarService: calendarService}
}

// GetMyEvents godoc
// @Summary      Получение событий календаря
// @Description  Возвращает события текущего пользователя за указанный период
// @Tags         calendar
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        from query string false "Начало периода (RFC3339)"
// @Param        to query string false "Конец периода (RFC3339)"
// @Success      200  {array}   calendar.EventResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /calendar/events/my [get]
func (h *CalendarHandler) GetMyEvents(c *gin.Context) {
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

	fromStr := c.DefaultQuery("from", time.Now().Format(time.RFC3339))
	toStr := c.DefaultQuery("to", time.Now().AddDate(0, 1, 0).Format(time.RFC3339))

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

	// Добавляем логирование
	log.Printf("[DEBUG] CalendarHandler.GetMyEvents: userID=%s, from=%s, to=%s", userID, from, to)

	events, err := h.calendarService.GetMyEvents(c.Request.Context(), userID, from, to)
	if err != nil {
		log.Printf("[DEBUG] Error in GetMyEvents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DEBUG] CalendarHandler.GetMyEvents: found %d events", len(events))
	c.JSON(http.StatusOK, events)
}

// CreateEvent godoc
// @Summary      Создание нового события
// @Description  Создаёт новое событие в календаре (для преподавателей/админов/держателей)
// @Tags         calendar
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body calendar.CreateEventRequest true "Данные события"
// @Success      201  {object}  calendar.EventResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /calendar/events [post]
func (h *CalendarHandler) CreateEvent(c *gin.Context) {
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

	var req calendar.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.calendarService.CreateEvent(c.Request.Context(), creatorID, &req)
	if err != nil {
		if err.Error() == "insufficient permissions to create events" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// UpdateEvent godoc
// @Summary      Обновление события
// @Description  Обновляет существующее событие
// @Tags         calendar
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Event ID"
// @Param        request body calendar.UpdateEventRequest true "Данные для обновления"
// @Success      200  {object}  calendar.EventResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /calendar/events/{id} [patch]
func (h *CalendarHandler) UpdateEvent(c *gin.Context) {
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

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req calendar.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.calendarService.UpdateEvent(c.Request.Context(), userID, eventID, &req)
	if err != nil {
		switch err.Error() {
		case "event not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions to modify this event":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent godoc
// @Summary      Удаление события
// @Description  Удаляет событие
// @Tags         calendar
// @Security     BearerAuth
// @Param        id path string true "Event ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /calendar/events/{id} [delete]
func (h *CalendarHandler) DeleteEvent(c *gin.Context) {
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

	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	err = h.calendarService.DeleteEvent(c.Request.Context(), userID, eventID)
	if err != nil {
		switch err.Error() {
		case "event not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions to delete this event":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
