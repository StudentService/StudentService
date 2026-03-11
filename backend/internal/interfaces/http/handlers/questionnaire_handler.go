package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/questionnaire"
)

type QuestionnaireHandler struct {
	questionnaireService *application.QuestionnaireService
}

func NewQuestionnaireHandler(questionnaireService *application.QuestionnaireService) *QuestionnaireHandler {
	return &QuestionnaireHandler{questionnaireService: questionnaireService}
}

// GetMyQuestionnaire godoc
// @Summary      Получение анкеты текущего пользователя
// @Description  Возвращает анкету текущего студента
// @Tags         questionnaire
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  questionnaire.QuestionnaireResponse
// @Success      204  "No Content - анкета не найдена"
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire/my [get]
func (h *QuestionnaireHandler) GetMyQuestionnaire(c *gin.Context) {
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

	questionnaire, err := h.questionnaireService.GetMyQuestionnaire(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if questionnaire == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, questionnaire)
}

// GetTemplate godoc
// @Summary      Получение шаблона анкеты
// @Description  Возвращает активный шаблон анкеты с описанием полей
// @Tags         questionnaire
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  questionnaire.TemplateResponse
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire/template [get]
func (h *QuestionnaireHandler) GetTemplate(c *gin.Context) {
	template, err := h.questionnaireService.GetTemplate(c.Request.Context())
	if err != nil {
		if err.Error() == "no active questionnaire template found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

// SubmitQuestionnaire godoc
// @Summary      Отправка анкеты
// @Description  Отправляет заполненную анкету на проверку
// @Tags         questionnaire
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body questionnaire.SubmitRequest true "Ответы на анкету"
// @Success      200  {object}  questionnaire.SubmitResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire/submit [post]
func (h *QuestionnaireHandler) SubmitQuestionnaire(c *gin.Context) {
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

	var req questionnaire.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.questionnaireService.SubmitQuestionnaire(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// SaveDraft godoc
// @Summary      Сохранение черновика
// @Description  Сохраняет черновик анкеты
// @Tags         questionnaire
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body questionnaire.SubmitRequest true "Ответы на анкету"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire/draft [post]
func (h *QuestionnaireHandler) SaveDraft(c *gin.Context) {
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

	var answers map[string]interface{} // меняем на map
	if err := c.ShouldBindJSON(&answers); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.questionnaireService.SaveDraft(c.Request.Context(), userID, answers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "draft saved successfully"})
}

// ReviewQuestionnaire godoc
// @Summary      Проверка анкеты (админ)
// @Description  Одобряет или отклоняет анкету студента
// @Tags         questionnaire
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "Questionnaire ID"
// @Param        request body questionnaire.ReviewRequest true "Решение по анкете"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire/{id}/review [post]
func (h *QuestionnaireHandler) ReviewQuestionnaire(c *gin.Context) {
	adminIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	adminID, err := uuid.Parse(adminIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	questionnaireID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid questionnaire id"})
		return
	}

	var req questionnaire.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.questionnaireService.ReviewQuestionnaire(c.Request.Context(), adminID, questionnaireID, &req); err != nil {
		switch err.Error() {
		case "questionnaire not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient permissions":
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "questionnaire reviewed successfully"})
}

// ListByStatus godoc
// @Summary      Список анкет по статусу (админ)
// @Description  Возвращает все анкеты с указанным статусом
// @Tags         questionnaire
// @Security     BearerAuth
// @Produce      json
// @Param        status query string true "Статус (draft, submitted, approved, rejected)"
// @Success      200  {array}   questionnaire.QuestionnaireResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questionnaire [get]
func (h *QuestionnaireHandler) ListByStatus(c *gin.Context) {
	// Проверяем права (только админ)
	role, exists := c.Get("user_role")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	statusStr := c.Query("status")
	if statusStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status parameter is required"})
		return
	}

	status := questionnaire.Status(statusStr)
	if status != questionnaire.StatusDraft &&
		status != questionnaire.StatusSubmitted &&
		status != questionnaire.StatusApproved &&
		status != questionnaire.StatusRejected {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	questionnaires, err := h.questionnaireService.ListQuestionnairesByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, questionnaires)
}
