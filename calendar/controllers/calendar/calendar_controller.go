package calendar

import (
	"github.com/atomic-blend/backend/calendar/repositories"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	"github.com/gin-gonic/gin"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles calendar operations
type Controller struct {
	calendarRepo repositories.CalendarRepositoryInterface
	amqpService  amqpinterfaces.AMQPServiceInterface
}

// NewCalendarController creates a new calendar controller
func NewCalendarController(calendarRepo repositories.CalendarRepositoryInterface, amqpService amqpinterfaces.AMQPServiceInterface) *Controller {
	return &Controller{calendarRepo: calendarRepo, amqpService: amqpService}
}

// SetupRoutes sets up calendar routes
func SetupRoutes(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	calendarRepo := repositories.NewCalendarRepository(database)
	calendarController := NewCalendarController(calendarRepo, amqpService)
	setupCalendarRoutes(router, calendarController)
}

// SetupRoutesWithMock sets up routes with mocks for testing
func SetupRoutesWithMock(router *gin.Engine, calendarRepo repositories.CalendarRepositoryInterface, amqpService amqpinterfaces.AMQPServiceInterface) {
	calendarController := NewCalendarController(calendarRepo, amqpService)
	setupCalendarRoutes(router, calendarController)
}

// setupCalendarRoutes registers routes
func setupCalendarRoutes(router *gin.Engine, calendarController *Controller) {
	calendarRoutes := router.Group("/calendar")
	auth.RequireAuth(calendarRoutes)
	{
		calendarRoutes.GET("", pagination.New(), calendarController.GetAllCalendars)
		calendarRoutes.GET("/since", pagination.New(), calendarController.GetCalendarsSince)
		calendarRoutes.GET(":id", calendarController.GetCalendarByID)
		calendarRoutes.POST("", calendarController.CreateCalendar)
		calendarRoutes.PUT(":id", calendarController.UpdateCalendar)
		calendarRoutes.DELETE(":id", calendarController.DeleteCalendar)
	}
}
