package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	taskRepo := repositories.NewTaskRepository(database)
	taskController := NewTaskController(taskRepo)
	setupTaskRoutes(router, taskController)
}

// SetupRoutesWithMock sets up the task routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, taskRepo repositories.TaskRepositoryInterface) {
	taskController := NewTaskController(taskRepo)
	setupTaskRoutes(router, taskController)
}

// setupTaskRoutes sets up the routes for task controller
func setupTaskRoutes(router *gin.Engine, taskController *TaskController) {
	taskRoutes := router.Group("/tasks")
	auth.RequireAuth(taskRoutes)
	{
		taskRoutes.GET("", taskController.GetAllTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.POST("", taskController.CreateTask)
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}
}
