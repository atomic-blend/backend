package tasks

import (
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/tests/mocks"
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

	// Create sample tags
	tagID1 := primitive.NewObjectID()
	tagID2 := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	name1 := "Tag 1"
	name2 := "Tag 2"

	tags := []*models.Tag{
		{
			ID:     &tagID1,
			UserID: &userID,
			Name:   name1,
		},
		{
			ID:     &tagID2,
			UserID: &userID,
			Name:   name2,
		},
	}

	return &models.TaskEntity{
		Title:       "Test Task",
		Description: &desc,
		Completed:   &completed,
		User:        userID,
		StartDate:   &now,
		EndDate:     &end,
		Reminders:   []*primitive.DateTime{&reminder1, &reminder2},
		Tags:        &tags,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}
}

func setupTest() (*gin.Engine, *mocks.MockTaskRepository, *mocks.MockTagRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTagRepo := new(mocks.MockTagRepository)
	taskController := NewTaskController(mockTaskRepo, mockTagRepo)

	// Set up routes with middleware
	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("", taskController.GetAllTasks)
		taskRoutes.GET("/:id", taskController.GetTaskByID)
		taskRoutes.POST("", taskController.CreateTask)
		taskRoutes.PUT("/:id", taskController.UpdateTask)
		taskRoutes.POST("/patch", taskController.Patch)
		taskRoutes.DELETE("/:id", taskController.DeleteTask)
	}

	return router, mockTaskRepo, mockTagRepo
}
