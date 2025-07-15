package notes

import (
	"productivity/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// UpdateNote updates an existing note
// @Summary Update note
// @Description Update an existing note
// @Tags Notes
// @Accept json
// @Produce json
// @Param id path string true "Note ID"
// @Param note body models.NoteEntity true "Note"
// @Success 200 {object} models.NoteEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes/{id} [put]
func (c *NoteController) UpdateNote(ctx *gin.Context) {
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to update this note"})
		return
	}

	var note models.NoteEntity
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set note owner to authenticated user (maintain ownership)
	note.User = authUser.UserID

	updatedNote, err := c.noteRepo.Update(ctx, id, &note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedNote)
}
