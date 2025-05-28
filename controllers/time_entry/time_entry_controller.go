package time_entry

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeEntryController struct {
	timeEntryRepository repositories.TimeEntryRepositoryInterface
}

func NewTimeEntryController(timeEntryRepository repositories.TimeEntryRepositoryInterface) *TimeEntryController {
	return &TimeEntryController{
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
