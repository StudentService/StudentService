package handlers

import (
	"backend/internal/domain/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe godoc
// @Summary      Получение профиля текущего пользователя
// @Description  Возвращает информацию о текущем авторизованном пользователе с группой и семестром
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  user.UserResponse
// @Failure      401  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// Используем новый метод с деталями
	u, err := h.userService.GetProfileWithDetails(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, u)
}

// GetUserByID с деталями TODO:в свагу добавить
func (h *UserHandler) GetUserByID(c *gin.Context) {
	currentUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	targetID := c.Param("id")
	if targetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id required"})
		return
	}

	// Получаем текущего пользователя для проверки прав
	currentUser, err := h.userService.GetProfile(c.Request.Context(), currentUserID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Получаем целевого пользователя с деталями
	targetUser, err := h.userService.GetProfileWithDetails(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if targetUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Проверяем права на просмотр
	// Для проверки прав нужно создать временный объект User
	tempUser := &user.User{
		ID:   targetUser.ID,
		Role: targetUser.Role,
	}
	if !currentUser.CanViewUser(tempUser) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to view this user"})
		return
	}

	c.JSON(http.StatusOK, targetUser)
}

// UpdateMe godoc
// @Summary      Обновление профиля текущего пользователя
// @Description  Обновляет информацию о текущем пользователе
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body object{first_name=string,last_name=string,username=string} true "Данные для обновления"
// @Success      200  {object}  user.User
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /users/me [patch]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var updateData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Username  string `json:"username"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущего пользователя
	u, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновляем поля
	if updateData.FirstName != "" {
		u.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		u.LastName = updateData.LastName
	}
	if updateData.Username != "" {
		u.Username = updateData.Username
	}

	// Сохраняем
	err = h.userService.UpdateProfile(c.Request.Context(), u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u.PasswordHash = ""
	c.JSON(http.StatusOK, u)
}
