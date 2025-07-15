package auth

import (
	"auth/models"
	"auth/repositories"
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
	userRepo          repositories.UserRepositoryInterface
	userRoleRepo      repositories.UserRoleRepositoryInterface
	resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface
}

// NewController creates a new auth controller
func NewController(userRepo repositories.UserRepositoryInterface, userRoleRepo repositories.UserRoleRepositoryInterface, resetPasswordRepo repositories.UserResetPasswordRequestRepositoryInterface) *Controller {
	return &Controller{
		userRepo:     userRepo,
		userRoleRepo: userRoleRepo,
		resetPasswordRepo: resetPasswordRepo,
	}
}
