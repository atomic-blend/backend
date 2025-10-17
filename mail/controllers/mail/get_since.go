package mail

import (
	"net/http"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
)

// GetMailsSince retrieves mails updated since a specific date for the authenticated user with pagination
// @Summary Get mails updated since date
// @Description Get mails updated since a specific date for the authenticated user with pagination
// @Tags Mail
// @Produce json
// @Param since query string true "Date in ISO8601 format (e.g., 2024-01-01T00:00:00Z)"
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} PaginatedMailResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/since [get]
func (c *Controller) GetMailsSince(ctx *gin.Context) {
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

	// Get pagination parameters from gin-pagination middleware
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// Get mails updated since the specified time with pagination
	mails, totalCount, err := c.mailRepo.GetSince(ctx, authUser.UserID, sinceTime, int64(page), int64(size))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (totalCount + int64(size) - 1) / int64(size)

	if mails == nil {
		mails = make([]*models.Mail, 0)
	}

	response := PaginatedMailResponse{
		Mails:      mails,
		TotalCount: totalCount,
		Page:       int64(page),
		Size:       int64(size),
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
