package tags

import (
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// TagController handles tag related operations
type TagController struct {
	tagRepo  repositories.TagRepositoryInterface
	taskRepo repositories.TaskRepositoryInterface
}

// NewTagController creates a new tag controller instance
func NewTagController(tagRepo repositories.TagRepositoryInterface, taskRepo repositories.TaskRepositoryInterface) *TagController {
	return &TagController{
		tagRepo:  tagRepo,
		taskRepo: taskRepo,
	}
}

// SetupRoutes sets up the routes for the tag controller
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	tagRepo := repositories.NewTagRepository(database)
	taskRepo := repositories.NewTaskRepository(database)
	tagController := NewTagController(tagRepo, taskRepo)
	setupTagRoutes(router, tagController)
}

// SetupRoutesWithMock sets up the tag routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, tagRepo repositories.TagRepositoryInterface, taskRepo repositories.TaskRepositoryInterface) {
	tagController := NewTagController(tagRepo, taskRepo)
	setupTagRoutes(router, tagController)
}

// setupTagRoutes sets up the routes for tag controller
func setupTagRoutes(router *gin.Engine, tagController *TagController) {
	tagRoutes := router.Group("/tags")
	auth.RequireAuth(tagRoutes)
	{
		tagRoutes.GET("", tagController.GetAllTags)
		tagRoutes.GET("/:id", tagController.GetTagByID)
		tagRoutes.POST("", tagController.CreateTag)
		tagRoutes.PUT("/:id", tagController.UpdateTag)
		tagRoutes.DELETE("/:id", tagController.DeleteTag)
	}
}
