package tasks

import (
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
)

// TaskController handles task related operations
type TaskController struct {
	taskRepo repositories.TaskRepositoryInterface
}

// NewTaskController creates a new task controller instance
func NewTaskController(taskRepo repositories.TaskRepositoryInterface) *TaskController {
	return &TaskController{
		taskRepo: taskRepo,
	}
}

// SetupRoutes sets up the task routes
func (c *TaskController) SetupRoutes(router *gin.RouterGroup) {
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", c.GetAllTasks)
		taskRoutes.GET("/:id", c.GetTaskByID)
		taskRoutes.POST("", c.CreateTask)
		taskRoutes.PUT("/:id", c.UpdateTask)
		taskRoutes.DELETE("/:id", c.DeleteTask)
	}
}
