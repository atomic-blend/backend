package sendmail

import (
	"net/http"

	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/models"

	"github.com/gin-gonic/gin"
)

// PaginatedSendMailResponse represents the paginated response for send mails
type PaginatedSendMailResponse struct {
	SendMails  []*models.SendMail `json:"send_mails"`
	TotalCount int64              `json:"total_count"`
	Page       int64              `json:"page,omitempty"`
	Size       int64              `json:"size,omitempty"`
	TotalPages int64              `json:"total_pages,omitempty"`
}

// GetAllSendMails retrieves all send mails for the authenticated user with pagination
// @Summary Get all send mails
// @Description Get all send mails for the authenticated user with pagination
// @Tags SendMail
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param size query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} PaginatedSendMailResponse
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /mail/send [get]
func (c *Controller) GetAllSendMails(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get pagination parameters from gin-pagination middleware
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	// Get send mails with pagination
	sendMails, totalCount, err := c.sendMailRepo.GetAll(ctx, authUser.UserID, int64(page), int64(size))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate total pages
	totalPages := (totalCount + int64(size) - 1) / int64(size)

	response := PaginatedSendMailResponse{
		SendMails:  sendMails,
		TotalCount: totalCount,
		Page:       int64(page),
		Size:       int64(size),
		TotalPages: totalPages,
	}

	ctx.JSON(http.StatusOK, response)
}
