package sendmail

import (
	"github.com/atomic-blend/backend/mail/repositories"
	userclient "github.com/atomic-blend/backend/shared/grpc/user"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	amqpinterfaces "github.com/atomic-blend/backend/shared/services/amqp/interfaces"
	s3service "github.com/atomic-blend/backend/shared/services/s3"
	s3interfaces "github.com/atomic-blend/backend/shared/services/s3/interfaces"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Controller handles send mail related operations
type Controller struct {
	sendMailRepo repositories.SendMailRepositoryInterface
	userClient   userclient.Interface
	amqpService  amqpinterfaces.AMQPServiceInterface
	s3Service    s3interfaces.S3ServiceInterface
}

// NewSendMailController creates a new send mail controller instance
func NewSendMailController(sendMailRepo repositories.SendMailRepositoryInterface, userClient userclient.Interface, amqpService amqpinterfaces.AMQPServiceInterface, s3Service s3interfaces.S3ServiceInterface) *Controller {
	return &Controller{
		sendMailRepo: sendMailRepo,
		userClient:   userClient,
		amqpService:  amqpService,
		s3Service:    s3Service,
	}
}

// SetupRoutes sets up the send mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	sendMailRepo := repositories.NewSendMailRepository(database)
	userClient, _ := userclient.NewUserClient()
	s3Service, _ := s3service.NewS3Service()
	sendMailController := NewSendMailController(sendMailRepo, userClient, amqpService, s3Service)
	setupSendMailRoutes(router, sendMailController)
}

// SetupRoutesWithMock sets up the send mail routes with mock services for testing
func SetupRoutesWithMock(router *gin.Engine, sendMailRepo repositories.SendMailRepositoryInterface, userClient userclient.Interface, amqpService amqpinterfaces.AMQPServiceInterface, s3Service s3interfaces.S3ServiceInterface) {
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
