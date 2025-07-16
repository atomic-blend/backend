package auth

import (
	"atomic-blend/backend/auth/utils/password"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// ConfirmResetPasswordRequest represents the request body for confirming a password reset
type ConfirmResetPasswordRequest struct {
	// ResetCode is the code sent to the user for resetting their password
	ResetCode string `json:"reset_code" binding:"required"`
	// ResetData indicates whether the user wants to reset their data
	ResetData *bool `json:"reset_data" binding:"required"`
	// NewPassword is the new password the user wants to set
	NewPassword string `json:"new_password" binding:"required"`

	// updated keyset for the user
	UserKey    string `json:"user_key" binding:"required"`
	UserSalt   string `json:"user_salt" binding:"required"`
	BackupKey  string `json:"backup_key" binding:"required"`
	BackupSalt string `json:"backup_salt" binding:"required"`
}

// ConfirmResetPassword handles the confirmation of a password reset
func (c *Controller) ConfirmResetPassword(ctx *gin.Context) {
	// Parse the request body
	var request ConfirmResetPasswordRequest
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

	// Update the user's password
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

	// Hash the new password
	hashedPassword, err := password.HashPassword(request.NewPassword)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Update the user's password in the database
	user.Password = &hashedPassword
	user.KeySet.UserKey = request.UserKey
	user.KeySet.Salt = request.UserSalt
	user.KeySet.BackupKey = request.BackupKey
	user.KeySet.MnemonicSalt = request.BackupSalt

	_, err = c.userRepo.Update(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// If ResetData is true, reset the user's data
	if *request.ResetData {
		err = c.userRepo.ResetAllUserData(ctx, *user.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to reset user data")
		}
	}

	// Delete the reset code from the database
	err = c.resetPasswordRepo.Delete(ctx, resetCode.UserID.Hex())
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete reset code")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete reset code"})
		return
	}

	// Respond with success
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	log.Info().Msg("Password reset successfully")
}
