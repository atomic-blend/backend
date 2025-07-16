package notes

import (
	"net/http"
	"atomic-blend/backend/productivity/auth"

	"github.com/gin-gonic/gin"
)

// GetAllNotes retrieves all notes for the authenticated user
// @Summary Get all notes
// @Description Get all notes for the authenticated user
// @Tags Notes
// @Produce json
// @Success 200 {array} models.NoteEntity
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes [get]
func (c *NoteController) GetAllNotes(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	notes, err := c.noteRepo.GetAll(ctx, &authUser.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notes)
}
