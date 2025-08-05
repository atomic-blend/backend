package mail

import (
	"net/http"
	"strconv"

	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/models"

	"github.com/gin-gonic/gin"
)

// PaginatedMailResponse represents the paginated response for mails
type PaginatedMailResponse struct {
	Mails      []*models.Mail `json:"mails"`
	TotalCount int64          `json:"total_count"`
	Page       int64          `json:"page,omitempty"`
	Limit      int64          `json:"limit,omitempty"`
	TotalPages int64          `json:"total_pages,omitempty"`
}

// GetAllMails retrieves all mails for the authenticated user with optional pagination
// @Summary Get all mails
// @Description Get all mails for the authenticated user with optional pagination
// @Tags Mail
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} PaginatedMailResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail [get]
func (c *MailController) GetAllMails(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	var mails []*models.Mail
	var totalCount int64
	var err error

	// Parse pagination parameters
	pageStr := ctx.DefaultQuery("page", "0")
	limitStr := ctx.DefaultQuery("limit", "0")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	// If both page and limit are zero, fetch all mails
	if page == 0 && limit == 0 {
		mails, totalCount, err = c.mailRepo.GetAll(ctx, authUser.UserID, 0, 0)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		response := PaginatedMailResponse{
			Mails:      mails,
			TotalCount: totalCount,
		}
		ctx.JSON(http.StatusOK, response)
		return
	}

	// If either page or limit is invalid, return 400
	if page <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}
	if limit <= 0 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter. Must be between 1 and 100"})
		return
	}

	// Get mails with pagination
	mails, totalCount, err = c.mailRepo.GetAll(ctx, authUser.UserID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (totalCount + limit - 1) / limit

	response := PaginatedMailResponse{
		Mails:      mails,
		TotalCount: totalCount,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
