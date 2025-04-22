package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/utils/password"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ConfirmResetPasswordRequest represents the request body for updating a password
type ConfirmResetPasswordRequest struct {
	Code         string `json:"code" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required"`
	UserKey      string `json:"user_key" binding:"required"`
	Salt         string `json:"salt" binding:"required"`
	MnemonicKey  string `json:"mnemonic_key" binding:"required"`
	MnemonicSalt string `json:"mnemonic_hash" binding:"required"`
}

// ConfirmResetPassword allows users to update their password
func (c *UserController) ConfirmResetPassword(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Parse and validate update request
	var updateReq ConfirmResetPasswordRequest
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

	// check if the reset code is valid
	if user.ResetPasswordCode == nil || *user.ResetPasswordCode != updateReq.Code {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid reset code"})
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
	user.KeySet.BackupKey = updateReq.MnemonicKey
	user.KeySet.MnemonicSalt = updateReq.MnemonicSalt

	// Update user in database
	if _, err := c.userRepo.Update(ctx, user); err != nil {
		log.Error().Err(err).Msg("Failed to update user password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password reseted successfully"})
}
