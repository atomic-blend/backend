package notes

import (
	"net/http"
	"strconv"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// PaginatedNoteResponse represents the paginated response for notes
type PaginatedNoteResponse struct {
	Notes      []*models.NoteEntity `json:"notes"`
	TotalCount int64                `json:"total_count"`
	Page       int64                `json:"page,omitempty"`
	Size       int64                `json:"size,omitempty"`
	TotalPages int64                `json:"total_pages,omitempty"`
}

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

	// Parse pagination parameters
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	var page, limit *int64
	var err error

	if pageStr != "" {
		pageVal, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || pageVal < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
		page = &pageVal
	}

	if limitStr != "" {
		limitVal, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limitVal < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		limit = &limitVal
	}

	notes, totalCount, err := c.noteRepo.GetAll(ctx, &authUser.UserID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pagination metadata
	var totalPages int64
	if limit != nil && *limit > 0 {
		totalPages = (totalCount + *limit - 1) / *limit
	}

	response := PaginatedNoteResponse{
		Notes:      notes,
		TotalCount: totalCount,
	}

	// Only include pagination metadata if pagination was used
	if page != nil && limit != nil && *page > 0 && *limit > 0 {
		response.Page = *page
		response.Size = *limit
		response.TotalPages = totalPages
	}

	ctx.JSON(http.StatusOK, response)
}
