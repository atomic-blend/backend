package auth

import (
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/utils/password"
	"github.com/atomic-blend/backend/auth/utils/resend"
	"bytes"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StartResetPasswordRequest represents the request body for starting a password reset
type StartResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// StartResetPassword handles the initiation of a password reset
func (c *Controller) StartResetPassword(ctx *gin.Context) {
	// Parse the request body
	var request StartResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := c.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		log.Error().Err(err).Msg("User not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// generate reset code
	resetCode, err := password.GenerateRandomString(8)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset code"})
		return
	}

	// template the html with gotemplate
	htmlTemplate, err := template.ParseFiles("./email_templates/reset_password/reset_password.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse HTML template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse HTML template"})
		return
	}

	textTemplate, err := template.ParseFiles("./email_templates/reset_password/reset_password.txt")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse text template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse text template"})
		return
	}

	// template the plain text with gotemplate
	var htmlContent bytes.Buffer
	err = htmlTemplate.Execute(&htmlContent, map[string]string{
		"code": resetCode,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to execute HTML template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute HTML template"})
		return
	}

	var textContent bytes.Buffer
	err = textTemplate.Execute(&textContent, map[string]string{
		"code": resetCode,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute text template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute text template"})
		return
	}

	// store the reset code in the database
	userResetPasswordRequest := &models.UserResetPassword{
		UserID:    user.ID,
		ResetCode: resetCode,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// check if a reset password request already exists for this user
	existingRequest, err := c.resetPasswordRepo.FindByUserID(ctx, user.ID.Hex())
	if err != nil {	
		log.Error().Err(err).Msg("Failed to find existing reset password request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find existing reset password request"})
		return
	}

	log.Info().Msgf("Existing request: %v", existingRequest)

	if existingRequest != nil {
		// delete the existing request
		err = c.resetPasswordRepo.Delete(ctx, existingRequest.UserID.Hex())
		if err != nil {
			log.Error().Err(err).Msg("Failed to delete existing reset password request")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing reset password request"})
			return
		}
	}

	//TODO store in db
	_, err = c.resetPasswordRepo.Create(ctx, userResetPasswordRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create reset password request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset password request"})
		return
	}

	// // send the email using the resend sdk
	emailClient := resend.NewResendClient(os.Getenv("RESEND_API_KEY"))
	sent, error := emailClient.Send(
		[]string{*user.BackupEmail},
		"Atomic Blend - Reset Password",
		htmlContent.String(),
		textContent.String(),
	)

	if error != nil {
		log.Error().Err(error).Msg("Failed to send email")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Reset password email sent successfully", "sent": sent})
}
