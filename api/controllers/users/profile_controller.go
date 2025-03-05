package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
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

// GetMyProfile returns the current authenticated user's profile
// @Summary Get authenticated user's profile
// @Description Returns the user profile for the currently authenticated user
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.UserEntity
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/profile [get]
func (c *ProfileController) GetMyProfile(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		// This shouldn't happen due to middleware, but just in case
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	// Fetch user profile from database
	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// For security, don't return the password hash
	user.Password = nil

	// Return user profile
	ctx.JSON(http.StatusOK, user)
}
