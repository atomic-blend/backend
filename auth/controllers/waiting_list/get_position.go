// Package waitinglist provides controllers for waiting list operations
package waitinglist

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWaitingListPositionRequest represents the structure for getting waiting list position request data
type GetWaitingListPositionRequest struct {
	Email         string `json:"email" binding:"required,email"`
	SecurityToken string `json:"securityToken" binding:"required"`
}

// GetWaitingListPosition returns the position and waiting list object for a given email and security token
func (c *Controller) GetWaitingListPosition(ctx *gin.Context) {
	var req GetWaitingListPositionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get the waiting list record by email
	waitingListRecord, err := c.waitingListRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_record"})
		return
	}

	if waitingListRecord == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "error_email_not_found_in_waiting_list"})
		return
	}

	// verify the security token matches
	if waitingListRecord.SecurityToken == nil || *waitingListRecord.SecurityToken != req.SecurityToken {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "error_invalid_security_token"})
		return
	}

	// get the total count of waiting list records
	totalCount, err := c.waitingListRepo.Count(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_count"})
		return
	}

	// get the position of this record (0-based index)
	position, err := c.waitingListRepo.GetPositionByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_position"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"entry":    waitingListRecord,
		"position": position,
		"total":    totalCount,
	})
}
