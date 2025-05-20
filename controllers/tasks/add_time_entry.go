package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AddTimeEntry adds a time entry to a task
// @Summary Add time entry
// @Description Add a time entry to a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param timeEntry body models.TimeEntry true "Time Entry"
// @Success 200 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id}/time-entries [post]
func (c *TaskController) AddTimeEntry(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	taskID := ctx.Param("id")
	if taskID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Get the task to verify ownership
	task, err := c.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching task: " + err.Error()})
		return
	}
	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Verify task ownership
	if task.User != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to modify this task"})
		return
	}

	// Parse time entry from request body
	var timeEntry models.TimeEntry
	if err := ctx.ShouldBindJSON(&timeEntry); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID if not provided
	if timeEntry.ID == nil {
		id := uuid.New().String()
		timeEntry.ID = &id
	}

	// Set created/updated timestamps
	now := time.Now().Format(time.RFC3339)
	timeEntry.CreatedAt = now
	timeEntry.UpdatedAt = now

	// Add time entry to task
	updatedTask, err := c.taskRepo.AddTimeEntry(ctx, taskID, &timeEntry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding time entry: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTask)
}
