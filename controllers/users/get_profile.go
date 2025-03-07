package users

import (
	"atomic_blend_api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetMyProfile fetches the profile of the authenticated user
func (c *UserController) GetMyProfile(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Fetch user from database using the ID from auth
	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve user profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Populate user roles
	err = c.userRoleRepo.PopulateRoles(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to populate user roles")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate user roles"})
		return
	}

	// Remove sensitive data before sending to client
	user.Password = nil

	// Return user profile
	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
