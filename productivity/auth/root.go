package auth

import (
	"github.com/gin-gonic/gin"
)

// RequireRoleMiddleware applies the auth middleware followed by role checking to a specific route group
// Example usage: RequireRoleMiddleware(router.Group("/admin"), "admin", userRepo)
func RequireRoleMiddleware(group *gin.RouterGroup, roleName string) *gin.RouterGroup {
	group.Use(Middleware())
	group.Use(requireRoleHandler(roleName))
	return group
}

// RequireAuth applies the auth middleware to a specific route group
// Example usage: RequireAuth(router.Group("/protected"))
func RequireAuth(group *gin.RouterGroup) *gin.RouterGroup {
	group.Use(Middleware())
	return group
}