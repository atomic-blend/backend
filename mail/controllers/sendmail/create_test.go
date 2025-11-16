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
	// Ensure PUBLIC_ADDRESS is set so controller generates Message-ID
	os.Setenv("PUBLIC_ADDRESS", "example.com")
	defer func() {
		os.Unsetenv("GO_ENV")
		os.Unsetenv("AWS_BUCKET")
		os.Unsetenv("PUBLIC_ADDRESS")
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
			name: "Success - generates Message-ID when missing",
			requestBody: models.RawMail{
				Headers: map[string]interface{}{
					"Subject": "Test Email",
					"From":    "test@example.com",
				},
				TextContent: "Test email content",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(arg interface{}) bool {
					sendMail, ok := arg.(*models.SendMail)
					if !ok {
						return false
					}

					if sendMail.Mail == nil || sendMail.SendStatus != models.SendStatusPending || sendMail.RetryCounter != nil {
						return false
					}

					if sendMail.Mail.Headers == nil {
						return false
					}

					// Validate Message-ID exists and is non-empty (headers are encrypted before repo call)
					if headers, ok := sendMail.Mail.Headers.(map[string]interface{}); ok {
						if v, exists := headers["Message-ID"]; exists {
							switch mv := v.(type) {
							case string:
								return mv != ""
							case []string:
								return len(mv) > 0 && mv[0] != ""
							}
						}
						return false
					}

					if headers2, ok := sendMail.Mail.Headers.(map[string][]string); ok {
						if mv, exists := headers2["Message-ID"]; exists {
							return len(mv) > 0 && mv[0] != ""
						}
						return false
					}

					return false
				})).Return(&models.SendMail{
					ID:           primitive.NewObjectID(),
					SendStatus:   models.SendStatusPending,
					RetryCounter: nil,
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
				mockAMQPService.On("PublishMessage", "mail", "sent", mock.MatchedBy(func(message map[string]interface{}) bool {
					_, hasSendMailID := message["send_mail_id"]
					_, hasContent := message["content"]
					return hasSendMailID && hasContent
				}), (*amqp.Table)(nil)).Return()
			},
			setupS3Mock: func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {
				mockS3Service.On("BulkUploadFiles", mock.Anything, mock.MatchedBy(func(payloads []*s3.PutObjectInput) bool {
					return len(payloads) == 0
				})).Return([]string{}, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Replaces Message-ID when provided",
			requestBody: models.RawMail{
				Headers: map[string]interface{}{
					"Message-ID": "<custom-id@example.com>",
					"Subject":    "Test Email",
					"From":       "test@example.com",
				},
				TextContent: "Test email content",
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(arg interface{}) bool {
					sendMail, ok := arg.(*models.SendMail)
					if !ok {
						return false
					}

					if sendMail.Mail == nil || sendMail.SendStatus != models.SendStatusPending {
						return false
					}

					if sendMail.Mail.Headers == nil {
						return false
					}

					// Message-ID should be present and encrypted (so just assert presence and non-empty)
					if headers, ok := sendMail.Mail.Headers.(map[string]interface{}); ok {
						if v, exists := headers["Message-ID"]; exists {
							switch mv := v.(type) {
							case string:
								return mv != "" && mv != "<custom-id@example.com>"
							case []string:
								return len(mv) > 0 && mv[0] != "" && mv[0] != "<custom-id@example.com>"
							}
						}
						return false
					}

					if headers2, ok := sendMail.Mail.Headers.(map[string][]string); ok {
						if arr, exists := headers2["Message-ID"]; exists {
							if len(arr) == 0 {
								return false
							}
							return arr[0] != "" && arr[0] != "<custom-id@example.com>"
						}
						return false
					}

					return false
				})).Return(&models.SendMail{
					ID:         primitive.NewObjectID(),
					SendStatus: models.SendStatusPending,
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
				mockAMQPService.On("PublishMessage", "mail", "sent", mock.MatchedBy(func(message map[string]interface{}) bool {
					_, hasSendMailID := message["send_mail_id"]
					_, hasContent := message["content"]
					return hasSendMailID && hasContent
				}), (*amqp.Table)(nil)).Return()
			},
			setupS3Mock: func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {
				mockS3Service.On("BulkUploadFiles", mock.Anything, mock.MatchedBy(func(payloads []*s3.PutObjectInput) bool {
					return len(payloads) == 0
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
