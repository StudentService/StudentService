package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/challenge"
)

type ChallengeHandler struct {
	challengeService *application.ChallengeService
}

func NewChallengeHandler(challengeService *application.ChallengeService) *ChallengeHandler {
	return &ChallengeHandler{challengeService: challengeService}
}

// GetMyChallenges godoc
// @Summary      Получение списка личных вызовов
// @Description  Возвращает все вызовы текущего пользователя
// @Tags         challenges
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}   challenge.ChallengeResponse
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /challenges/my [get]
func (h *ChallengeHandler) GetMyChallenges(c *gin.Context) {
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

	challenges, err := h.challengeService.GetMyChallenges(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, challenges)
}

// GetChallengeByID godoc
// @Summary      Получение вызова по ID
// @Description  Возвращает детальную информацию о конкретном вызове
// @Tags         challenges
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Challenge ID"
// @Success      200  {object}  challenge.ChallengeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /challenges/{id} [get]
func (h *ChallengeHandler) GetChallengeByID(c *gin.Context) {
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

	challengeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	challenge, err := h.challengeService.GetChallengeByID(c.Request.Context(), userID, challengeID)
	if err != nil {
		switch err.Error() {
		case "challenge not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "access denied":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, challenge)
}

// CreateChallenge godoc
// @Summary      Создание нового вызова
// @Description  Создает новый личный вызов для студента
// @Tags         challenges
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body challenge.CreateChallengeRequest true "Данные вызова"
// @Success      201  {object}  challenge.ChallengeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /challenges [post]
func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
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

	var req challenge.CreateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge, err := h.challengeService.CreateChallenge(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

// UpdateChallenge godoc
// @Summary      Обновление вызова
// @Description  Обновляет существующий вызов
// @Tags         challenges
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Challenge ID"
// @Param        request body challenge.UpdateChallengeRequest true "Данные для обновления"
// @Success      200  {object}  challenge.ChallengeResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /challenges/{id} [patch]
func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
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

	challengeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	var req challenge.UpdateChallengeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	challenge, err := h.challengeService.UpdateChallenge(c.Request.Context(), userID, challengeID, &req)
	if err != nil {
		switch err.Error() {
		case "challenge not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "access denied":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, challenge)
}

// DeleteChallenge godoc
// @Summary      Удаление вызова
// @Description  Удаляет вызов (только если он в статусе draft)
// @Tags         challenges
// @Security     BearerAuth
// @Param        id path string true "Challenge ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /challenges/{id} [delete]
func (h *ChallengeHandler) DeleteChallenge(c *gin.Context) {
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

	challengeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	err = h.challengeService.DeleteChallenge(c.Request.Context(), userID, challengeID)
	if err != nil {
		switch err.Error() {
		case "challenge not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "access denied":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
