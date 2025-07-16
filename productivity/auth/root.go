package auth

import (
	"atomic-blend/backend/productivity/repositories"
	"atomic-blend/backend/productivity/utils/db"

	"github.com/gin-gonic/gin"
)

// RequireRoleMiddleware applies the auth middleware followed by role checking to a specific route group
// Example usage: RequireRoleMiddleware(router.Group("/admin"), "admin", userRepo)
func RequireRoleMiddleware(group *gin.RouterGroup, roleName string) *gin.RouterGroup {
	userRepo := repositories.NewUserRepository(db.Database)
	userRoleRepo := repositories.NewUserRoleRepository(db.Database)
	group.Use(Middleware())
	group.Use(requireRoleHandler(roleName, userRepo, userRoleRepo))
	return group
}

// RequireAuth applies the auth middleware to a specific route group
// Example usage: RequireAuth(router.Group("/protected"))
func RequireAuth(group *gin.RouterGroup) *gin.RouterGroup {
	group.Use(Middleware())
	return group
}