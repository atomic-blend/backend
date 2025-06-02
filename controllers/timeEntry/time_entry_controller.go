package timeentrycontroller

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles time entry related operations
type Controller struct {
	timeEntryRepository repositories.TimeEntryRepositoryInterface
}

// NewTimeEntryController creates a new instance of TimeEntryController
func NewTimeEntryController(timeEntryRepository repositories.TimeEntryRepositoryInterface) *Controller {
	return &Controller{
		timeEntryRepository: timeEntryRepository,
	}
}

// SetupRoutes configures all time entry related routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	timeEntryRepo := repositories.NewTimeEntryRepository(database)
	timeEntryController := NewTimeEntryController(timeEntryRepo)

	// Apply authentication middleware to all time entry routes
	timeEntryGroup := router.Group("/time-entries")
	timeEntryGroup.Use(auth.Middleware())
	{
		timeEntryGroup.GET("", timeEntryController.GetAll)
		timeEntryGroup.GET("/:id", timeEntryController.GetByID)
		timeEntryGroup.POST("", timeEntryController.Create)
		timeEntryGroup.PUT("/:id", timeEntryController.Update)
		timeEntryGroup.DELETE("/:id", timeEntryController.Delete)
	}
}
