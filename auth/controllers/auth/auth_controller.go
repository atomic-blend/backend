package auth

import (
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/shared/models"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRequest represents the structure for registration request data
type RegisterRequest struct {
	Email    string                `json:"email" binding:"required,email"`
	KeySet   *models.EncryptionKey `json:"keySet" binding:"required"`
	Password string                `json:"password" binding:"required,min=8"` // Minimum 8 characters
}

// Response represents the structure for authentication response data
type Response struct {
	User         *models.UserEntity `json:"user"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	ExpiresAt    int64              `json:"expiresAt"`
}

// Controller handles auth-related operations
type Controller struct {
	userRepo          userrepo.UserRepositoryInterface
	userRoleRepo      userrolerepo.UserRoleRepositoryInterface
	resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface
}

// NewController creates a new auth controller
func NewController(userRepo userrepo.UserRepositoryInterface, userRoleRepo userrolerepo.UserRoleRepositoryInterface, resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface) *Controller {
	return &Controller{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		resetPasswordRepo: resetPasswordRepo,
	}
}


// SetupRoutes configures all auth-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(database)
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/refresh", authController.RefreshToken)
		authGroup.POST("/reset-password", authController.StartResetPassword)
		authGroup.POST("/reset-password/backup-key", authController.GetBackupKeyForResetPassword)
		authGroup.POST("/reset-password/confirm", authController.ConfirmResetPassword)
	}
}