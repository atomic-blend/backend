package mail

import (
	"net/http"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// CleanupTrash cleans up trash mails for the authenticated user
func (c *Controller) CleanupTrash(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	log.Debug().Str("user_id", authUser.UserID.Hex()).Msg("Starting trash cleanup for user")

	// Call the repository method with the authenticated user's ID
	err := c.mailRepo.CleanupTrash(ctx, &authUser.UserID)
	if err != nil {
		log.Error().Err(err).Str("user_id", authUser.UserID.Hex()).Msg("Failed to cleanup trash")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cleanup trash"})
		return
	}

	log.Info().Str("user_id", authUser.UserID.Hex()).Msg("Trash cleanup completed successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "Trash cleanup completed successfully"})
}
