package controllers

import (
	"github.com/atomic-blend/backend/mail/controllers/mail"
	"github.com/atomic-blend/backend/mail/controllers/sendmail"
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAllControllers sets up all controllers
func SetupAllControllers(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	// Setup mail controller
	mail.SetupRoutes(router, database)
	sendmail.SetupRoutes(router, database, amqpService)
}
