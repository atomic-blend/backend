package notes

import (
	"atomic_blend_api/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetNoteByID retrieves a note by its ID
// @Summary Get note by ID
// @Description Get a note by its ID
// @Tags Notes
// @Produce json
// @Param id path string true "Note ID"
// @Success 200 {object} models.NoteEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes/{id} [get]
func (c *NoteController) GetNoteByID(ctx *gin.Context) {
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

	note, err := c.noteRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if note == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	// Check if the note belongs to the authenticated user
	if note.User != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this note"})
		return
	}

	ctx.JSON(http.StatusOK, note)
}
