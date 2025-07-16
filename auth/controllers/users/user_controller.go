package users

import (
	"github.com/atomic-blend/backend/auth/auth"
	"github.com/atomic-blend/backend/auth/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserController handles user profile related operations
type UserController struct {
	userRepo     repositories.UserRepositoryInterface
	userRoleRepo repositories.UserRoleRepositoryInterface
}

// NewUserController creates a new profile controller instance
func NewUserController(userRepo repositories.UserRepositoryInterface, userRoleRepo repositories.UserRoleRepositoryInterface) *UserController {
	return &UserController{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := repositories.NewUserRepository(database)
	userRoleRepo := repositories.NewUserRoleRepository(database)
	userController := NewUserController(userRepo, userRoleRepo)

	// Public user routes (if any)
	userGroup := router.Group("/users")

	// Protected user routes (require authentication)
	protectedUserRoutes := auth.RequireAuth(userGroup)
	{
		protectedUserRoutes.GET("/profile", userController.GetMyProfile)
		protectedUserRoutes.PUT("/profile", userController.UpdateProfile)
		protectedUserRoutes.PUT("/password", userController.UpdatePassword)
		protectedUserRoutes.DELETE("/me", userController.DeleteAccount)
		protectedUserRoutes.PUT("/device", userController.UpdateDeviceInfo)
		// Add more protected routes here as needed
	}
}
