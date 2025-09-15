package mail

import (
	"github.com/atomic-blend/backend/mail/repositories"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles mail related operations
type Controller struct {
	mailRepo repositories.MailRepositoryInterface
}

// NewMailController creates a new mail controller instance
func NewMailController(mailRepo repositories.MailRepositoryInterface) *Controller {
	return &Controller{
		mailRepo: mailRepo,
	}
}

// SetupRoutes sets up the mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	mailRepo := repositories.NewMailRepository(database)
	mailController := NewMailController(mailRepo)
	setupMailRoutes(router, mailController)
}

// SetupRoutesWithMock sets up the mail routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, mailRepo repositories.MailRepositoryInterface) {
	mailController := NewMailController(mailRepo)
	setupMailRoutes(router, mailController)
}

// setupMailRoutes sets up the routes for mail controller
func setupMailRoutes(router *gin.Engine, mailController *Controller) {
	mailRoutes := router.Group("/mail")
	auth.RequireAuth(mailRoutes)
	{
		mailRoutes.GET("/", pagination.New(), mailController.GetAllMails)
		mailRoutes.GET("/:id", mailController.GetMailByID)
		mailRoutes.PUT("/actions", mailController.PutMailActions)
		mailRoutes.POST("/trash/empty", mailController.CleanupTrash)
	}
}
