package folder

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles folder related operations
type Controller struct {
	folderRepo repositories.FolderRepositoryInterface
}

// NewFolderController creates a new folder controller instance
func NewFolderController(folderRepo repositories.FolderRepositoryInterface) *Controller {
	return &Controller{
		folderRepo: folderRepo,
	}
}

// SetupRoutes sets up the routes for the folder controller
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	folderRepo := repositories.NewFolderRepository(database)
	folderController := NewFolderController(folderRepo)
	setupFolderRoutes(router, folderController)
}

// SetupRoutesWithMock sets up the folder routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, folderRepo repositories.FolderRepositoryInterface) {
	folderController := NewFolderController(folderRepo)
	setupFolderRoutes(router, folderController)
}

// setupFolderRoutes sets up the routes for folder controller
func setupFolderRoutes(router *gin.Engine, folderController *Controller) {
	folderRoutes := router.Group("/folders")
	auth.RequireAuth(folderRoutes)
	{
		// Folder endpoints
		folderRoutes.POST("", folderController.CreateFolder)
		folderRoutes.GET("", folderController.GetAllFolders)
		folderRoutes.PUT("/:id", folderController.UpdateFolder)
		folderRoutes.DELETE("/:id", folderController.DeleteFolder)
	}
}
