// Package config provides configuration-related HTTP handlers for the auth service.
package config

import (
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles config related operations
type Controller struct {
	userRepo user.Interface
}

// NewConfigController creates a new config controller instance
func NewConfigController(userRepo user.Interface) *Controller {
	return &Controller{
		userRepo: userRepo,
	}
}

// SetupRoutes sets up the config routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	userRepo := user.NewUserRepository(database)
	configController := NewConfigController(userRepo)
	setupConfigRoutes(router, configController)
}

// SetupRoutesWithMock sets up the config routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, database *mongo.Database) {
	userRepo := user.NewUserRepository(database)
	configController := NewConfigController(userRepo)
	setupConfigRoutes(router, configController)
}

// setupConfigRoutes sets up the routes for config controller
func setupConfigRoutes(router *gin.Engine, configController *Controller) {
	configRoutes := router.Group("/config")
	{
		configRoutes.GET("", configController.GetConfig)
	}
}
