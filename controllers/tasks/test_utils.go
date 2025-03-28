package tasks

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestTask() *models.TaskEntity {
	desc := "Test Description"
	completed := false
	now := primitive.NewDateTimeFromTime(time.Now())
	end := primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour))
	reminder1 := primitive.NewDateTimeFromTime(time.Now().Add(12 * time.Hour))
	reminder2 := primitive.NewDateTimeFromTime(time.Now().Add(18 * time.Hour))

	return &models.TaskEntity{
		Title:       "Test Task",
		Description: &desc,
		Completed:   &completed,
		User:        primitive.NewObjectID(),
		StartDate:   &now,
		EndDate:     &end,
		Reminders:   []*primitive.DateTime{&reminder1, &reminder2},
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
}

func setupTest() (*gin.Engine, *mocks.MockTaskRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockTaskRepository)
	taskController := NewTaskController(mockRepo)

	// Set up routes with middleware
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", taskController.GetAllTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.POST("", taskController.CreateTask)
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}

	return router, mockRepo
}
