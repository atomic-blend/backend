package notes

import (
	"net/http"
	"strconv"
	"time"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// GetNotesSince retrieves notes updated since a specific date for the authenticated user with pagination
// @Summary Get notes updated since date
// @Description Get notes updated since a specific date for the authenticated user with pagination
// @Tags Notes
// @Produce json
// @Param since query string true "Date in ISO8601 format (e.g., 2024-01-01T00:00:00Z)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} PaginatedNoteResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /notes/since [get]
func (c *NoteController) GetNotesSince(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get the 'since' query parameter
	sinceStr := ctx.Query("since")
	if sinceStr == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: since"})
		return
	}

	// Parse the ISO8601 date string
	sinceTime, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		// Try parsing with timezone offset format
		sinceTime, err = time.Parse("2006-01-02T15:04:05-07:00", sinceStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected ISO8601 format (e.g., 2024-01-01T00:00:00Z or 2024-01-01T12:30:45+02:00)"})
			return
		}
	}

	// Parse pagination parameters
	pageStr := ctx.Query("page")
	limitStr := ctx.Query("limit")

	var page, limit *int64

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

	// Get notes updated since the specified time with pagination
	notes, totalCount, err := c.noteRepo.GetSince(ctx, authUser.UserID, sinceTime, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	var totalPages int64
	if limit != nil && *limit > 0 {
		totalPages = (totalCount + *limit - 1) / *limit
	}

	if notes == nil {
		notes = make([]*models.NoteEntity, 0)
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
