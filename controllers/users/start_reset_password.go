package users

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/utils/password"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *UserController) startResetPassword(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// generate reset code
	resetCode, err := password.GenerateRandomString(4)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset code"})
		return
	}

	// send email to account email

	// store the reset code in the database
}