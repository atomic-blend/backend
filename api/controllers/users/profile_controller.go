package users

import (
	"atomic_blend_api/repositories"
)

// ProfileController handles user profile related operations
type ProfileController struct {
	userRepo *repositories.UserRepository
	userRoleRepo *repositories.UserRoleRepository
}

// NewProfileController creates a new profile controller instance
func NewProfileController(userRepo *repositories.UserRepository, userRoleRepo *repositories.UserRoleRepository) *ProfileController {
	return &ProfileController{
		userRepo: userRepo,
		userRoleRepo: userRoleRepo,
	}
}