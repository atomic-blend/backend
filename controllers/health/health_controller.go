package health

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthController struct {
	database *mongo.Database
}

type HealthResponse struct {
	Up bool `json:"up"`
	DB bool `json:"db"`
}

func NewHealthController(database *mongo.Database) *HealthController {
	return &HealthController{
		database: database,
	}
}

func (c *HealthController) GetHealth(ctx *gin.Context) {
	// Check database connectivity
	dbStatus := true
	client := c.database.Client()
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctxTimeout, readpref.Primary())
	if err != nil {
		dbStatus = false
	}

	response := HealthResponse{
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
