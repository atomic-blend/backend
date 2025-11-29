package calendar

import (
	"github.com/atomic-blend/backend/calendar/repositories"
	authmw "github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouterWithMockRepoForTests(mockRepo repositories.CalendarRepositoryInterface) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// inject auth user middleware for tests using the mock repo owner
	testUserID := primitive.NewObjectID()
	r.Use(func(c *gin.Context) {
		c.Set("authUser", &authmw.UserAuthInfo{UserID: testUserID})
		c.Next()
	})

	ctrl := NewCalendarController(mockRepo, nil)
	grp := r.Group("/calendar")
	{
		grp.GET("", func(c *gin.Context) { c.Set("page", 1); c.Set("size", 10); ctrl.GetAllCalendars(c) })
		grp.GET("/since", func(c *gin.Context) { c.Set("page", 1); c.Set("size", 10); ctrl.GetCalendarsSince(c) })
		grp.GET(":id", ctrl.GetCalendarByID)
		grp.POST("", ctrl.CreateCalendar)
		grp.PUT(":id", ctrl.UpdateCalendar)
		grp.DELETE(":id", ctrl.DeleteCalendar)
	}
	return r
}
