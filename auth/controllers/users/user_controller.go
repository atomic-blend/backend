package users

import (
	productivityclient "github.com/atomic-blend/backend/shared/grpc/productivity"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserController handles user profile related operations
type UserController struct {
	userRepo           userrepo.Interface
	userRoleRepo       userrolerepo.Interface
	productivityClient productivityclient.Interface
}

// NewUserController creates a new profile controller instance
func NewUserController(userRepo userrepo.Interface, userRoleRepo userrolerepo.Interface, productivityClient productivityclient.Interface) *UserController {
	return &UserController{
		userRepo:           userRepo,
		userRoleRepo:       userRoleRepo,
		productivityClient: productivityClient,
	}
}

// SetupRoutes configures all user-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)

	// // Create productivity client
	productivityClient, err := productivityclient.NewProductivityClient()
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
