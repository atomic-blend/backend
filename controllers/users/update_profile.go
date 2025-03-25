package users

import (
	"atomic_blend_api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UpdateProfileRequest represents the data needed to update a user profile
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
	// Add more fields here as needed, such as name, phone, preferences, etc.
}

// UpdateProfile allows users to update their profile information
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse and validate update request
	var updateReq UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Fetch the current user data
	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve user profile for update")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Only update fields that were provided in the request
	if updateReq.Email != "" {
		// Check if email is already in use by another user
		if existingUser, err := c.userRepo.FindByEmail(ctx, updateReq.Email); err == nil && existingUser != nil {
			if existingUser.ID.Hex() != user.ID.Hex() {
				ctx.JSON(http.StatusConflict, gin.H{"error": "Email is already in use"})
				return
			}
		}

		emailVal := updateReq.Email
		user.Email = &emailVal
	}

	// Add handling for additional fields here

	// Update user in database
	updatedUser, err := c.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user profile")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Remove sensitive data before sending response
	updatedUser.Password = nil
	updatedUser.KeySet = nil

	// Return updated user profile
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"data":    updatedUser,
	})
}
