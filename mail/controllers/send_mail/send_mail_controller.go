package send_mail

import (
	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/grpc/clients"
	"github.com/atomic-blend/backend/mail/grpc/interfaces"
	mailinterfaces "github.com/atomic-blend/backend/mail/services/interfaces"
	"github.com/atomic-blend/backend/mail/repositories"
	"github.com/atomic-blend/backend/mail/services"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles send mail related operations
type Controller struct {
	sendMailRepo repositories.SendMailRepositoryInterface
	userClient   interfaces.UserClientInterface
	amqpService  mailinterfaces.AMQPServiceInterface
	s3Service    mailinterfaces.S3ServiceInterface
}

// NewSendMailController creates a new send mail controller instance
func NewSendMailController(sendMailRepo repositories.SendMailRepositoryInterface, userClient interfaces.UserClientInterface, amqpService mailinterfaces.AMQPServiceInterface, s3Service mailinterfaces.S3ServiceInterface) *Controller {
	return &Controller{
		sendMailRepo: sendMailRepo,
		userClient:   userClient,
		amqpService:  amqpService,
		s3Service:    s3Service,
	}
}

// SetupRoutes sets up the send mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database) {
	sendMailRepo := repositories.NewSendMailRepository(database)
	userClient, _ := clients.NewUserClient()
	amqpService := services.NewAMQPService()
	s3Service, _ := services.NewS3Service()
	sendMailController := NewSendMailController(sendMailRepo, userClient, amqpService, s3Service)
	setupSendMailRoutes(router, sendMailController)
}

// SetupRoutesWithMock sets up the send mail routes with mock services for testing
func SetupRoutesWithMock(router *gin.Engine, sendMailRepo repositories.SendMailRepositoryInterface, userClient interfaces.UserClientInterface, amqpService mailinterfaces.AMQPServiceInterface, s3Service mailinterfaces.S3ServiceInterface) {
	sendMailController := NewSendMailController(sendMailRepo, userClient, amqpService, s3Service)
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
