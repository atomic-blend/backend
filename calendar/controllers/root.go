package controllers

import (
	calendar "github.com/atomic-blend/backend/calendar/controllers/calendar"
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAllControllers sets up all controllers
func SetupAllControllers(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	// calendar controller
	calendar.SetupRoutes(router, database, amqpService)
}
