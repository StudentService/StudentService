package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"backend/internal/application"
	"backend/internal/domain/group"
)

type GroupHandler struct {
	groupService *application.GroupService
}

func NewGroupHandler(groupService *application.GroupService) *GroupHandler {
	return &GroupHandler{groupService: groupService}
}

// GetAllGroups godoc
// @Summary      Получение всех групп
// @Tags         groups
// @Security     BearerAuth
// @Produce      json
// @Success      200  {array}  group.Group
// @Router       /groups [get]
func (h *GroupHandler) GetAllGroups(c *gin.Context) {
	groups, err := h.groupService.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// GetGroupByID godoc
// @Summary      Получение группы по ID
// @Tags         groups
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Success      200  {object}  group.Group
// @Router       /groups/{id} [get]
func (h *GroupHandler) GetGroupByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	g, err := h.groupService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if g == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, g)
}

// CreateGroup godoc
// @Summary      Создание группы
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body group.CreateGroupRequest true "Данные группы"
// @Success      201  {object}  group.Group
// @Router       /groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req group.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.groupService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

// UpdateGroup godoc
// @Summary      Обновление группы
// @Tags         groups
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Param        request body group.UpdateGroupRequest true "Данные для обновления"
// @Success      200  {object}  group.Group
// @Router       /groups/{id} [patch]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	var req group.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	g, err := h.groupService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

// DeleteGroup godoc
// @Summary      Удаление группы
// @Tags         groups
// @Security     BearerAuth
// @Param        id path string true "Group ID"
// @Success      204  "No Content"
// @Router       /groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}

	if err := h.groupService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
