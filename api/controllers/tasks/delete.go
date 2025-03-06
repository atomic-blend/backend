package tasks

import (
	"net/http"

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
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [delete]
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id := ctx.Param("id")

	// Check if task exists
	existingTask, err := c.taskRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingTask == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Delete the task
	err = c.taskRepo.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
