package send_mail

import (
	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/repositories"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles send mail related operations
type Controller struct {
	sendMailRepo repositories.SendMailRepositoryInterface
}

// NewSendMailController creates a new send mail controller instance
func NewSendMailController(sendMailRepo repositories.SendMailRepositoryInterface) *Controller {
	return &Controller{
		sendMailRepo: sendMailRepo,
	}
}

// SetupRoutes sets up the send mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	sendMailRepo := repositories.NewSendMailRepository(database)
	sendMailController := NewSendMailController(sendMailRepo)
	setupSendMailRoutes(router, sendMailController)
}

// SetupRoutesWithMock sets up the send mail routes with a mock repository for testing
func SetupRoutesWithMock(router *gin.Engine, sendMailRepo repositories.SendMailRepositoryInterface) {
	sendMailController := NewSendMailController(sendMailRepo)
	setupSendMailRoutes(router, sendMailController)
}

// setupSendMailRoutes sets up the routes for send mail controller
func setupSendMailRoutes(router *gin.Engine, sendMailController *Controller) {
	sendMailRoutes := router.Group("/mail/send")
	auth.RequireAuth(sendMailRoutes)
	{
		sendMailRoutes.GET("/", pagination.New(), sendMailController.GetAllSendMails)
		sendMailRoutes.GET("/:id", sendMailController.GetSendMailByID)
		sendMailRoutes.POST("/", sendMailController.CreateSendMail)
		sendMailRoutes.DELETE("/:id", sendMailController.DeleteSendMail)
	}
}
