package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllTasks retrieves all tasks
// @Summary Get all tasks
// @Description Get all tasks, optionally filtered by user
// @Tags Tasks
// @Produce json
// @Param user_id query string false "Filter by user ID"
// @Success 200 {array} models.TaskEntity
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [get]
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	// Get user ID from query param if provided
	userIDStr := ctx.Query("user_id")
	var userID *primitive.ObjectID

	// If user_id is provided, convert it to ObjectID
	if userIDStr != "" {
		objID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}
		userID = &objID
	}

	tasks, err := c.taskRepo.GetAll(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}
