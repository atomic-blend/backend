package auth

import (
	"atomic_blend_api/controllers/auth"
	"atomic_blend_api/repositories"
	"atomic_blend_api/utils/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configures all auth-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := repositories.NewUserRepository(database)
	userRoleRepo := repositories.NewUserRoleRepository(database)
	authController := auth.NewController(userRepo, userRoleRepo)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh", authController.RefreshToken)
	}
}

// RequireRoleMiddleware applies the auth middleware followed by role checking to a specific route group
// Example usage: RequireRoleMiddleware(router.Group("/admin"), "admin", userRepo)
func RequireRole(group *gin.RouterGroup, roleName string) *gin.RouterGroup {
	userRepo := repositories.NewUserRepository(db.Database)
	userRoleRepo := repositories.NewUserRoleRepository(db.Database)
	group.Use(AuthMiddleware())
	group.Use(requireRoleHandler(roleName, userRepo, userRoleRepo))
	return group
}

// RequireAuth applies the auth middleware to a specific route group
// Example usage: RequireAuth(router.Group("/protected"))
func RequireAuth(group *gin.RouterGroup) *gin.RouterGroup {
	group.Use(AuthMiddleware())
	return group
}
