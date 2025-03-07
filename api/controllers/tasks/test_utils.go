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

	return &models.TaskEntity{
		Title:       "Test Task",
		Description: &desc,
		Completed:   &completed,
		StartDate:   &now,
		EndDate:     &end,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
}

func setupTest() (*gin.Engine, *mocks.MockTaskRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockRepo := new(mocks.MockTaskRepository)

	// Use the new SetupRoutesWithMock function instead of SetupRoutes
	SetupRoutesWithMock(router, mockRepo)

	return router, mockRepo
}