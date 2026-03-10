package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
	"backend/internal/domain/auth"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      Авторизация пользователя
// @Description  Получение JWT токенов по email и паролю
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.LoginRequest true "Данные для входа"
// @Success      200  {object}  auth.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Обновление токенов
// @Description  Получение новой пары токенов по refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RefreshRequest true "Refresh token"
// @Success      200  {object}  auth.AuthResponse
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req auth.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Logout godoc
// @Summary      Выход из системы
// @Description  Инвалидация refresh token (клиент должен удалить токены)
// @Tags         auth
// @Success      200  {object}  map[string]string
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// В простейшем случае просто возвращаем успех
	// В реальном проекте здесь можно инвалидировать refresh token
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// Register godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя и возвращает данные профиля + токены
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body auth.RegisterRequest true "Данные для регистрации"
// @Success      201  {object}  auth.RegisterResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp) // Теперь resp.User уже без PasswordHash
}
