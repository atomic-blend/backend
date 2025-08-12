package notes

import (
	"net/http"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

// DeleteNote deletes a note by its ID
// @Summary Delete note
// @Description Delete a note by its ID
// @Tags Notes
// @Param id path string true "Note ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes/{id} [delete]
func (c *NoteController) DeleteNote(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	id := ctx.Param("id")
	if id == "" || strings.TrimSpace(id) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Note ID is required"})
		return
	}

	// Check if the note exists and belongs to the user
	existingNote, err := c.noteRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingNote == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Check if the note belongs to the authenticated user
	if existingNote.User != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this note"})
		return
	}

	err = c.noteRepo.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
