package mail

import (
	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// MailController handles mail related operations
type MailController struct {
	mailRepo repositories.MailRepositoryInterface
}

// NewMailController creates a new mail controller instance
func NewMailController(mailRepo repositories.MailRepositoryInterface) *MailController {
	return &MailController{
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
func setupMailRoutes(router *gin.Engine, mailController *MailController) {
	mailRoutes := router.Group("/mail")
	auth.RequireAuth(mailRoutes)
	{
		mailRoutes.GET("", mailController.GetAllMails)
		mailRoutes.GET("/:id", mailController.GetMailByID)
	}
}
