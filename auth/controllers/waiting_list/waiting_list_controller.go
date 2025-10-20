package waitinglist

import (
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles waiting list-related operations
type Controller struct {
	waitingListRepo repositories.WaitingListRepositoryInterface
}

// NewController creates a new waiting list controller
func NewController(waitingListRepo repositories.WaitingListRepositoryInterface) *Controller {
	return &Controller{waitingListRepo: waitingListRepo}
}

// SetupRoutes configures the waiting list routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	waitingListRepo := repositories.NewWaitingListRepository(database)
	waitingListController := NewController(waitingListRepo)
	waitingListGroup := router.Group("/waiting-list")
	{
		waitingListGroup.POST("", waitingListController.JoinWaitingList)
	}
}