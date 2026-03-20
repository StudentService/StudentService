package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"backend/internal/application"
)

// RequirePermission проверяет права доступа по роли
func RequirePermission(rbacService *application.RBACService, resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		allowed := rbacService.CheckUserPermission(role.(string), resource, action)
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied: insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin проверяет, что пользователь - админ
func RequireAdmin(rbacService *application.RBACService) gin.HandlerFunc {
	return RequirePermission(rbacService, "*", "*")
}

// RequireTeacherOrAdmin проверяет, что пользователь - преподаватель или админ
func RequireTeacherOrAdmin(rbacService *application.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		if role == "admin" || role == "teacher" {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: teacher or admin required"})
		c.Abort()
	}
}

// RequireHolderOrAdmin проверяет, что пользователь - держатель или админ
func RequireHolderOrAdmin(rbacService *application.RBACService) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			c.Abort()
			return
		}

		if role == "admin" || role == "holder" {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: holder or admin required"})
		c.Abort()
	}
}
