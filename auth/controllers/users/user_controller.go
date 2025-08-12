package users

import (
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/auth/grpc/clients"
	"github.com/atomic-blend/backend/auth/grpc/interfaces"
	"github.com/atomic-blend/backend/auth/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserController handles user profile related operations
type UserController struct {
	userRepo           repositories.UserRepositoryInterface
	userRoleRepo       repositories.UserRoleRepositoryInterface
	productivityClient interfaces.ProductivityClientInterface
}

// NewUserController creates a new profile controller instance
func NewUserController(userRepo repositories.UserRepositoryInterface, userRoleRepo repositories.UserRoleRepositoryInterface, productivityClient interfaces.ProductivityClientInterface) *UserController {
	return &UserController{
		userRepo:           userRepo,
		userRoleRepo:       userRoleRepo,
		productivityClient: productivityClient,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := repositories.NewUserRepository(database)
	userRoleRepo := repositories.NewUserRoleRepository(database)

	// Create productivity client
	productivityClient, err := clients.NewProductivityClient()
	if err != nil {
		panic("Failed to create productivity client: " + err.Error())
	}

	userController := NewUserController(userRepo, userRoleRepo, productivityClient)

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
