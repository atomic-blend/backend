package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/utils/password"
	"atomic_blend_api/utils/resend"
	"bytes"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (c *Controller) StartResetPassword(ctx *gin.Context) {
	// Get authenticated user from context
	// authUser := auth.GetAuthUser(ctx)
	// if authUser == nil {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
	// 	return
	// }

	// Parse the request body
	var request ResetPasswordRequest
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
	userResetPasswordRequest := &models.UserResetPasswordRequest{
		UserID:    user.ID,
		ResetCode: resetCode,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
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
		[]string{*user.Email},
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
