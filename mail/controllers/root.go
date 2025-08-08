package controllers

import (
	"github.com/atomic-blend/backend/mail/controllers/mail"
	"github.com/atomic-blend/backend/mail/controllers/send_mail"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupAllControllers sets up all controllers
func SetupAllControllers(router *gin.Engine, database *mongo.Database) {
	// Setup mail controller
	mail.SetupRoutes(router, database)
	send_mail.SetupRoutes(router, database)
}
