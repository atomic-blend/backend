package waitinglist

import (
	"github.com/atomic-blend/backend/auth/repositories"
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles waiting list-related operations
type Controller struct {
	waitingListRepo repositories.WaitingListRepositoryInterface
	amqpService     amqpinterfaces.AMQPServiceInterface
}

// NewController creates a new waiting list controller
func NewController(waitingListRepo repositories.WaitingListRepositoryInterface, amqpService amqpinterfaces.AMQPServiceInterface) *Controller {
	return &Controller{waitingListRepo: waitingListRepo, amqpService: amqpService}
}

// SetupRoutes configures the waiting list routes
func SetupRoutes(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	waitingListRepo := repositories.NewWaitingListRepository(database)
	waitingListController := NewController(waitingListRepo, amqpService)
	waitingListGroup := router.Group("/auth/waiting-list")
	{
		waitingListGroup.POST("", waitingListController.JoinWaitingList)
		waitingListGroup.POST("/position", waitingListController.GetWaitingListPosition)
	}
}
