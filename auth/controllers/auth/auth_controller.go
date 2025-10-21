// Package auth provides authentication and authorization
package auth

import (
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/grpc/gen/mailserver/v1/mailserverv1connect"
	mailserver "github.com/atomic-blend/backend/shared/grpc/mail-server"
	"github.com/atomic-blend/backend/shared/models"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterRequest represents the structure for registration request data
type RegisterRequest struct {
	Email           string                `json:"email" binding:"required,email"`
	WaitingListCode *string               `json:"waitingListCode"`
	BackupEmail     *string               `json:"backupEmail"`
	FirstName       *string               `json:"firstName"`
	LastName        *string               `json:"lastName"`
	KeySet          *models.EncryptionKey `json:"keySet" binding:"required"`
	Password        string                `json:"password" binding:"required,min=8"` // Minimum 8 characters
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
	userRepo          userrepo.Interface
	userRoleRepo      userrolerepo.Interface
	resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface
	waitingListRepo repositories.WaitingListRepositoryInterface
	mailServerClient  mailserverv1connect.MailServerServiceClient
}

// NewController creates a new auth controller
func NewController(userRepo userrepo.Interface, userRoleRepo userrolerepo.Interface, resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface, waitingListRepo repositories.WaitingListRepositoryInterface, mailServerClient mailserverv1connect.MailServerServiceClient) *Controller {
	return &Controller{
		userRepo:          userRepo,
		userRoleRepo:      userRoleRepo,
		resetPasswordRepo: resetPasswordRepo,
		waitingListRepo:   waitingListRepo,
		mailServerClient:  mailServerClient,
	}
}

// SetupRoutes configures all auth-related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(database)
	mailServerClient, _ := mailserver.NewMailServerClient()
	waitingListRepo := repositories.NewWaitingListRepository(database)
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo, waitingListRepo, mailServerClient)

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
