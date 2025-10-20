package waitinglist

import (
	"github.com/atomic-blend/backend/auth/repositories"
	mailserver "github.com/atomic-blend/backend/shared/grpc/mail-server"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles waiting list-related operations
type Controller struct {
	waitingListRepo  repositories.WaitingListRepositoryInterface
	mailServerClient mailserver.Interface
}

// NewController creates a new waiting list controller
func NewController(waitingListRepo repositories.WaitingListRepositoryInterface, mailServerClient mailserver.Interface) *Controller {
	return &Controller{waitingListRepo: waitingListRepo, mailServerClient: mailServerClient}
}

// SetupRoutes configures the waiting list routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	waitingListRepo := repositories.NewWaitingListRepository(database)
	mailServerClient, _ := mailserver.NewMailServerClient()
	waitingListController := NewController(waitingListRepo, mailServerClient)
	waitingListGroup := router.Group("/auth/waiting-list")
	{
		waitingListGroup.POST("", waitingListController.JoinWaitingList)
		waitingListGroup.POST("/position", waitingListController.GetWaitingListPosition)
	}
}
