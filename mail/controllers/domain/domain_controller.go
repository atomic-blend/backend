package mail

import (
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles mail related operations
type Controller struct {
}

// NewMailController creates a new mail controller instance
func NewMailController() *Controller {
	return &Controller{}
}

// SetupRoutes sets up the mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	mailController := NewMailController()
	setupMailRoutes(router, mailController)
}

// SetupRoutesWithMock sets up the mail routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine) {
	mailController := NewMailController()
	setupMailRoutes(router, mailController)
}

// setupMailRoutes sets up the routes for mail controller
func setupMailRoutes(router *gin.Engine, mailController *Controller) {
	mailRoutes := router.Group("/domain")
	auth.RequireAuth(mailRoutes)
	{
	}
}
