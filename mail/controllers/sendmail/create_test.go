package sendmail

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"connectrpc.com/connect"
	userv1 "github.com/atomic-blend/backend/grpc/gen/user/v1"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	amqpservice "github.com/atomic-blend/backend/shared/services/amqp"
	s3service "github.com/atomic-blend/backend/shared/services/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/streadway/amqp"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendMailController_CreateSendMail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// Set test environment
	os.Setenv("GO_ENV", "test")
	os.Setenv("AWS_BUCKET", "test-bucket")
	defer func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("AWS_BUCKET")
	}()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		setupMock      func(*mocks.MockSendMailRepository, primitive.ObjectID)
		setupUserMock  func(*mocks.MockUserClient, primitive.ObjectID)
		setupAMQPMock  func(*amqpservice.MockAMQPService, primitive.ObjectID)
		setupS3Mock    func(*s3service.MockS3Service, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name: "Success",
			requestBody: models.RawMail{
				Headers: map[string]interface{}{
					"Subject": "Test Email",
					"From":    "test@example.com",
				},
				TextContent: "Test email content",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(sendMail *models.SendMail) bool {
					return sendMail.Mail != nil &&
						sendMail.SendStatus == models.SendStatusPending &&
						sendMail.RetryCounter == nil // Should be nil when created
				})).Return(&models.SendMail{
					ID:           primitive.NewObjectID(),
					SendStatus:   models.SendStatusPending,
					RetryCounter: nil, // Should be nil when created
				}, nil)
			},
			setupUserMock: func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {
				mockUserClient.On("GetUserPublicKey", mock.Anything, mock.MatchedBy(func(req *connect.Request[userv1.GetUserPublicKeyRequest]) bool {
					return req.Msg.Id == userID.Hex()
				})).Return(&connect.Response[userv1.GetUserPublicKeyResponse]{
					Msg: &userv1.GetUserPublicKeyResponse{
						PublicKey: "age1jl76v4rmz5ukg9danl3v0zmyet9sqejmngs52wj9m497wgd02s9quq4qfl",
						UserId:    userID.Hex(),
					},
				}, nil)
			},
			setupAMQPMock: func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {
				// Mock AMQP publish message
				mockAMQPService.On("PublishMessage", "mail", "sent", mock.MatchedBy(func(message map[string]interface{}) bool {
					// Verify that the message contains the expected fields
					_, hasSendMailID := message["send_mail_id"]
					_, hasContent := message["content"]
					return hasSendMailID && hasContent
				}), (*amqp.Table)(nil)).Return()
			},
			setupS3Mock: func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {
				// Mock S3 operations for attachments (empty attachments in this test)
				mockS3Service.On("BulkUploadFiles", mock.Anything, mock.MatchedBy(func(payloads []*s3.PutObjectInput) bool {
					return len(payloads) == 0 // No attachments expected
				})).Return([]string{}, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Missing mail field",
			requestBody:    models.RawMail{},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unauthorized",
			requestBody: models.RawMail{
				TextContent: "Test",
			},
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth:      func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			mockUserClient := &mocks.MockUserClient{}
			mockAMQPService := &amqpservice.MockAMQPService{}
			mockS3Service := &s3service.MockS3Service{}
			userID := primitive.NewObjectID()

			tt.setupMock(mockRepo, userID)
			tt.setupUserMock(mockUserClient, userID)
			tt.setupAMQPMock(mockAMQPService, userID)
			tt.setupS3Mock(mockS3Service, userID)

			controller := NewSendMailController(mockRepo, mockUserClient, mockAMQPService, mockS3Service)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			sendMailRoutes := router.Group("/mail/send")
			{
				sendMailRoutes.POST("", controller.CreateSendMail)
			}

			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req, _ := http.NewRequest("POST", "/mail/send", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
			mockUserClient.AssertExpectations(t)
			mockAMQPService.AssertExpectations(t)
			mockS3Service.AssertExpectations(t)
		})
	}
}
