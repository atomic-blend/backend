package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// GetBackupKeyForResetPasswordRequest represents the request body for getting the backup key for resetting a password
type GetBackupKeyForResetPasswordRequest struct {
	// ResetCode is the code sent to the user for resetting their password
	ResetCode string `json:"reset_code" binding:"required"`
}

// GetBackupKeyForResetPassword handles the retrieval of the backup key for resetting a password
func (c *Controller) GetBackupKeyForResetPassword(ctx *gin.Context) {
	// Parse the request body
	var request GetBackupKeyForResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate the reset code
	resetCode, err := c.resetPasswordRepo.FindByResetCode(ctx, request.ResetCode)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find reset code")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find reset code"})
		return
	}
	if resetCode == nil {
		log.Error().Msg("Reset code not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reset code not found"})
		return
	}

	// Check if the reset code is expired
	if resetCode.CreatedAt.Time().Add(5 * time.Minute).Before(time.Now()) {
		log.Error().Msg("Reset code expired")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Reset code expired"})
		return
	}

	// get the user associated with the reset code
	user, err := c.userRepo.FindByID(ctx, *resetCode.UserID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find user"})
		return
	}
	if user == nil {
		log.Error().Msg("User not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"backup_key":   user.KeySet.BackupKey,
		"backup_salt":  user.KeySet.MnemonicSalt,
	})
}