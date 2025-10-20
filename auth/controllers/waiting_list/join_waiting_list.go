package waitinglist

import (
	"net/http"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// JoinWaitingListRequest represents the structure for join waiting list request data
type JoinWaitingListRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// JoinWaitingList creates a new waiting list record and returns a success message
func (c *Controller) JoinWaitingList(ctx *gin.Context) {
	var req JoinWaitingListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if the email is already in the waiting list
	waitingListRecord, err := c.waitingListRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_record"})
		return
	}

	if waitingListRecord != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "error_email_already_in_waiting_list"})
		return
	}

	// get the number of waiting list records
	waitingBeforeCount, err := c.waitingListRepo.Count(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_getting_waiting_list_count"})
		return
	}

	// create a new waiting list record
	now := primitive.NewDateTimeFromTime(time.Now())
	waitingListRecord, err = c.waitingListRepo.Create(ctx, &waitinglist.WaitingList{
		Email:     req.Email,
		CreatedAt: &now,
		UpdatedAt: &now,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error_creating_waiting_list_record"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"entry":        waitingListRecord,
		"before_count": waitingBeforeCount,
	})
}
