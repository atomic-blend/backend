package users

import (
	"atomic_blend_api/repositories"
)

// ProfileController handles user profile related operations
type ProfileController struct {
	userRepo *repositories.UserRepository
}

// NewProfileController creates a new profile controller instance
func NewProfileController(userRepo *repositories.UserRepository) *ProfileController {
	return &ProfileController{
		userRepo: userRepo,
	}
}