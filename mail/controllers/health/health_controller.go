package health

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Controller is a controller for health check actions
type Controller struct {
	database *mongo.Database
}

// Response is a health check response
type Response struct {
	Up bool `json:"up"`
	DB bool `json:"db"`
}

// NewHealthController creates a new health controller
func NewHealthController(database *mongo.Database) *Controller {
	return &Controller{
		database: database,
	}
}

// GetHealth returns the health status of the API
func (c *Controller) GetHealth(ctx *gin.Context) {
	// Check database connectivity
	dbStatus := true
	client := c.database.Client()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctxTimeout, readpref.Primary())
	if err != nil {
		dbStatus = false
	}

	response := Response{
		Up: true, // API is up if this code is executing
		DB: dbStatus,
	}

	ctx.JSON(200, response)
}

// SetupRoutes registers health check routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	healthController := NewHealthController(database)

	// Public endpoint - no authentication required
	router.GET("/health", healthController.GetHealth)
}
