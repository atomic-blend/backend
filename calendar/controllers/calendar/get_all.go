package calendar

import (
	"net/http"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// PaginatedCalendarResponse represents paginated calendars
type PaginatedCalendarResponse struct {
	Calendars  []*models.Calendar `json:"calendars"`
	TotalCount int64              `json:"total_count"`
	Page       int64              `json:"page,omitempty"`
	Size       int64              `json:"size,omitempty"`
	TotalPages int64              `json:"total_pages"`
}

// GetAllCalendars returns all calendars for the user with pagination
func (c *Controller) GetAllCalendars(ctx *gin.Context) {
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	calendars, totalCount, err := c.calendarRepo.GetAll(ctx, authUser.UserID, int64(page), int64(size))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalCount + int64(size) - 1) / int64(size)

	response := PaginatedCalendarResponse{
		Calendars:  calendars,
		TotalCount: totalCount,
		Page:       int64(page),
		Size:       int64(size),
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
