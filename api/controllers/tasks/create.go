package tasks

import (
	"atomic_blend_api/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateTask creates a new task
// @Summary Create task
// @Description Create a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.TaskEntity true "Task"
// @Success 201 {object} models.TaskEntity
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /tasks [post]
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task models.TaskEntity

	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values if needed
	if task.Completed == nil {
		completed := false
		task.Completed = &completed
	}

	createdTask, err := c.taskRepo.Create(ctx, &task)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdTask)
}
