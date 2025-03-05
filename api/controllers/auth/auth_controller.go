package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
)

// RegisterRequest represents the structure for registration request data
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"` // Minimum 8 characters
}

// AuthResponse represents the structure for authentication response data
type AuthResponse struct {
	User         *models.UserEntity `json:"user"`
	AccessToken  string             `json:"accessToken"`
	RefreshToken string             `json:"refreshToken"`
	ExpiresAt    int64              `json:"expiresAt"`
}

// Controller handles auth-related operations
type Controller struct {
	userRepo *repositories.UserRepository
}

// NewController creates a new auth controller
func NewController(userRepo *repositories.UserRepository) *Controller {
	return &Controller{
		userRepo: userRepo,
	}
}

