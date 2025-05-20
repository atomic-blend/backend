package tasks

import (
	"atomic_blend_api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RemoveTimeEntry removes a time entry from a task
// @Summary Remove time entry
// @Description Remove a time entry from a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param entryId path string true "Time Entry ID"
// @Success 200 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id}/time-entries/{entryId} [delete]
func (c *TaskController) RemoveTimeEntry(ctx *gin.Context) {
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

	timeEntryID := ctx.Param("entryID")
	if timeEntryID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Time Entry ID is required"})
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

	// Remove time entry from task
	updatedTask, err := c.taskRepo.RemoveTimeEntry(ctx, taskID, timeEntryID)
	if err != nil {
		if err.Error() == "no time entries found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Time entry not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error removing time entry: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedTask)
}
