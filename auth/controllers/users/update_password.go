package users

import (
	"atomic-blend/backend/auth/auth"
	"atomic-blend/backend/auth/utils/password"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// UpdatePasswordRequest represents the request body for updating a password
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	UserKey     string `json:"user_key" binding:"required"`
	Salt        string `json:"salt" binding:"required"`
}

// UpdatePassword allows users to update their password
func (c *UserController) UpdatePassword(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse and validate update request
	var updateReq UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&updateReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Fetch the current user data
	user, err := c.userRepo.FindByID(ctx, authUser.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve user profile for password update")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user profile"})
		return
	}

	// Verify the old password
	if !password.CheckPassword(updateReq.OldPassword, *user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	// Encrypt the new password
	newPasswordHash, err := password.HashPassword(updateReq.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to encrypt new password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt new password"})
		return
	}

	user.Password = &newPasswordHash
	user.KeySet.UserKey = updateReq.UserKey
	user.KeySet.Salt = updateReq.Salt

	// Update user in database
	if _, err := c.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
