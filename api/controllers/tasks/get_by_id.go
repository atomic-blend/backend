package tasks

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaskByID retrieves a task by its ID
// @Summary Get task by ID
// @Description Get a task by its ID
// @Tags Tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks/{id} [get]
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	id := ctx.Param("id")

	task, err := c.taskRepo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if task == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}
