// Package config provides configuration-related HTTP handlers for the auth service.
package config

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles config related operations
type Controller struct {
}

// NewConfigController creates a new config controller instance
func NewConfigController() *Controller {
	return &Controller{}
}

// SetupRoutes sets up the config routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	configController := NewConfigController()
	setupConfigRoutes(router, configController)
}

// SetupRoutesWithMock sets up the config routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine) {
	configController := NewConfigController()
	setupConfigRoutes(router, configController)
}

// setupConfigRoutes sets up the routes for config controller
func setupConfigRoutes(router *gin.Engine, configController *Controller) {
	configRoutes := router.Group("/config")
	{
		configRoutes.GET("", configController.AvailableAccountDomain)
	}
}
