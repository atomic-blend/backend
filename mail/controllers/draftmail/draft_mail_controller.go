package draftmail

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

// Controller handles draft mail related operations
type Controller struct {
	draftMailRepo repositories.DraftMailRepositoryInterface
	userClient    userclient.Interface
	amqpService   amqpinterfaces.AMQPServiceInterface
	s3Service     s3interfaces.S3ServiceInterface
}

// NewDraftMailController creates a new draft mail controller instance
func NewDraftMailController(draftMailRepo repositories.DraftMailRepositoryInterface, userClient userclient.Interface, amqpService amqpinterfaces.AMQPServiceInterface, s3Service s3interfaces.S3ServiceInterface) *Controller {
	return &Controller{
		draftMailRepo: draftMailRepo,
		userClient:    userClient,
		amqpService:   amqpService,
		s3Service:     s3Service,
	}
}

// SetupRoutes sets up the draft mail routes
func SetupRoutes(router *gin.Engine, database *mongo.Database, amqpService amqpinterfaces.AMQPServiceInterface) {
	draftMailRepo := repositories.NewDraftMailRepository(database)
	userClient, _ := userclient.NewUserClient()
	s3Service, _ := s3service.NewS3Service()
	draftMailController := NewDraftMailController(draftMailRepo, userClient, amqpService, s3Service)
	setupDraftMailRoutes(router, draftMailController)
}

// SetupRoutesWithMock sets up the draft mail routes with mock services for testing
func SetupRoutesWithMock(router *gin.Engine, draftMailRepo repositories.DraftMailRepositoryInterface, userClient userclient.Interface, amqpService amqpinterfaces.AMQPServiceInterface, s3Service s3interfaces.S3ServiceInterface) {
	draftMailController := NewDraftMailController(draftMailRepo, userClient, amqpService, s3Service)
	setupDraftMailRoutes(router, draftMailController)
}

// setupDraftMailRoutes sets up the routes for draft mail controller
func setupDraftMailRoutes(router *gin.Engine, draftMailController *Controller) {
	draftMailRoutes := router.Group("/mail/draft")
	auth.RequireAuth(draftMailRoutes)
	{
		draftMailRoutes.GET("", pagination.New(), draftMailController.GetAllDraftMails)
		draftMailRoutes.GET("/:id", draftMailController.GetDraftMailByID)
		draftMailRoutes.POST("", draftMailController.CreateDraftMail)
		draftMailRoutes.PUT("/:id", draftMailController.UpdateDraftMail)
		draftMailRoutes.DELETE("/:id", draftMailController.DeleteDraftMail)
	}
}
