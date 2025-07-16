package tasks

import (
	"net/http"
	"atomic-blend/backend/productivity/auth"

	"github.com/gin-gonic/gin"
)

// DeleteTask deletes a task by its ID
// @Summary Delete task
// @Description Delete a task by its ID
// @Tags Tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [delete]
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	// Get authenticated user from context
	authUser := auth.GetAuthUser(ctx)
	if authUser == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}

	// Get task ID from URL
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// First get the task to check ownership
	task, err := c.taskRepo.GetByID(ctx, id)
	// Added nil check
	if err != nil || task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Check if the authenticated user owns this task
	if task.User != authUser.UserID {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this task"})
		return
	}

	err = c.taskRepo.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
