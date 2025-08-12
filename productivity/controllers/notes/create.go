package notes

import (
	"net/http"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/productivity/models"

	"github.com/gin-gonic/gin"
)

// CreateNote creates a new note
// @Summary Create note
// @Description Create a new note
// @Tags Notes
// @Accept json
// @Produce json
// @Param note body models.NoteEntity true "Note"
// @Success 201 {object} models.NoteEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes [post]
func (c *NoteController) CreateNote(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var note models.NoteEntity
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set note owner to authenticated user
	note.User = authUser.UserID

	createdNote, err := c.noteRepo.Create(ctx, &note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdNote)
}
