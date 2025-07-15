package tags

import (
	"productivity/models"
	"productivity/tests/mocks"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTestTag() *models.Tag {
	name := "Test Tag"
	color := "#FF5733"
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Tag{
		Name:      name,
		Color:     &color,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func setupTest() (*gin.Engine, *mocks.MockTagRepository, *mocks.MockTaskRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockTagRepo := new(mocks.MockTagRepository)
	mockTaskRepo := new(mocks.MockTaskRepository)
	tagController := NewTagController(mockTagRepo, mockTaskRepo)

	// Set up routes with middleware
	tagRoutes := router.Group("/tags")
	{
		tagRoutes.GET("", tagController.GetAllTags)
		tagRoutes.GET("/:id", tagController.GetTagByID)
		tagRoutes.POST("", tagController.CreateTag)
		tagRoutes.PUT("/:id", tagController.UpdateTag)
		tagRoutes.DELETE("/:id", tagController.DeleteTag)
	}

	return router, mockTagRepo, mockTaskRepo
}
