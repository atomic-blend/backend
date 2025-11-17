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

	originalMailIDHex := "507f1f77bcf86cd799439011"
	originalMailID, _ := primitive.ObjectIDFromHex(originalMailIDHex)

	tests := []struct {
		name           string
		requestBody    func() interface{}
		expectedStatus int
		setupMock      func(*mocks.MockSendMailRepository, primitive.ObjectID)
		setupUserMock  func(*mocks.MockUserClient, primitive.ObjectID)
		setupAMQPMock  func(*amqpservice.MockAMQPService, primitive.ObjectID)
		setupS3Mock    func(*s3service.MockS3Service, primitive.ObjectID)
		setupMailMock  func(*mocks.MockMailRepository, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name: "Success - generates Message-ID when missing",
			requestBody: func() interface{} {
				return models.RawMail{
					Headers: map[string]interface{}{
						"Subject": "Test Email",
						"From":    "test@example.com",
					},
					TextContent: "Test email content",
				}
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
					if v, exists := sendMail.Mail.Headers["Message-ID"]; exists {
						switch mv := v.(type) {
						case string:
							return mv != ""
						case []string:
							return len(mv) > 0 && mv[0] != ""
						case []interface{}:
							if len(mv) > 0 {
								if s, ok := mv[0].(string); ok {
									return s != ""
								}
							}
						}
					}

					// Headers surface as map[string]interface{}; values can still be string or []string

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
			setupMailMock: func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Replaces Message-ID when provided",
			requestBody: func() interface{} {
				return models.RawMail{
					Headers: map[string]interface{}{
						"Message-ID": "<custom-id@example.com>",
						"Subject":    "Test Email",
						"From":       "test@example.com",
					},
					TextContent: "Test email content",
				}
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
					if v, exists := sendMail.Mail.Headers["Message-ID"]; exists {
						switch mv := v.(type) {
						case string:
							return mv != "" && mv != "<custom-id@example.com>"
						case []string:
							return len(mv) > 0 && mv[0] != "" && mv[0] != "<custom-id@example.com>"
						case []interface{}:
							if len(mv) > 0 {
								if s, ok := mv[0].(string); ok {
									return s != "" && s != "<custom-id@example.com>"
								}
							}
						}
					}

					// No extra assertion here: handled above via the map[string]interface{} branch

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
			setupMailMock: func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success - reply to email sets In-Reply-To and References",
			requestBody: func() interface{} {
				return models.RawMail{
					Headers: map[string]interface{}{
						"Subject": "Re: Test Email",
						"From":    "test@example.com",
					},
					TextContent: "Reply content",
					InReplyTo:   &originalMailIDHex,
				}
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

					// Check Message-ID is present
					if v, exists := sendMail.Mail.Headers["Message-ID"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check In-Reply-To is present
					if v, exists := sendMail.Mail.Headers["In-Reply-To"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check References is present
					if v, exists := sendMail.Mail.Headers["References"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check InReplyTo is set in the Mail entity
					if sendMail.Mail.InReplyTo == nil || *sendMail.Mail.InReplyTo != originalMailID {
						return false
					}

					return true
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
			setupMailMock: func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mockMailRepo.On("GetByID", mock.Anything, originalMailID).Return(&models.Mail{
					Headers: map[string]interface{}{
						"Message-ID": "<original@example.com>",
					},
				}, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Success - reply to email with existing References",
			requestBody: func() interface{} {
				return models.RawMail{
					Headers: map[string]interface{}{
						"Subject": "Re: Test Email",
						"From":    "test@example.com",
					},
					TextContent: "Reply content",
					InReplyTo:   &originalMailIDHex,
				}
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

					// Check Message-ID is present
					if v, exists := sendMail.Mail.Headers["Message-ID"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check In-Reply-To is present
					if v, exists := sendMail.Mail.Headers["In-Reply-To"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check References is present
					if v, exists := sendMail.Mail.Headers["References"]; !exists {
						return false
					} else {
						switch mv := v.(type) {
						case string:
							if mv == "" {
								return false
							}
						case []string:
							if len(mv) == 0 || mv[0] == "" {
								return false
							}
						case []interface{}:
							if len(mv) == 0 {
								if s, ok := mv[0].(string); !ok || s == "" {
									return false
								}
							}
						}
					}

					// Check InReplyTo is set in the Mail entity
					if sendMail.Mail.InReplyTo == nil || *sendMail.Mail.InReplyTo != originalMailID {
						return false
					}

					return true
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
			setupMailMock: func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mockMailRepo.On("GetByID", mock.Anything, originalMailID).Return(&models.Mail{
					Headers: map[string]interface{}{
						"Message-ID":  "<original@example.com>",
						"References": "<ref1@example.com> <ref2@example.com>",
					},
				}, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid request body",
			requestBody:    func() interface{} { return "invalid json" },
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupMailMock:  func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Missing mail field",
			requestBody:    func() interface{} { return models.RawMail{} },
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupMailMock:  func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unauthorized",
			requestBody: func() interface{} {
				return models.RawMail{
					TextContent: "Test",
				}
			},
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupUserMock:  func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupAMQPMock:  func(mockAMQPService *amqpservice.MockAMQPService, userID primitive.ObjectID) {},
			setupS3Mock:    func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupMailMock:  func(mockMailRepo *mocks.MockMailRepository, userID primitive.ObjectID) {},
			setupAuth:      func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			mockUserClient := &mocks.MockUserClient{}
			mockAMQPService := &amqpservice.MockAMQPService{}
			mockS3Service := &s3service.MockS3Service{}
			mockMailRepo := &mocks.MockMailRepository{}
			userID := primitive.NewObjectID()

			tt.setupMock(mockRepo, userID)
			tt.setupUserMock(mockUserClient, userID)
			tt.setupAMQPMock(mockAMQPService, userID)
			tt.setupS3Mock(mockS3Service, userID)
			tt.setupMailMock(mockMailRepo, userID)

			controller := NewSendMailController(mockRepo, mockMailRepo, mockUserClient, mockAMQPService, mockS3Service)

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
			json.NewEncoder(&body).Encode(tt.requestBody())

			req, _ := http.NewRequest("POST", "/mail/send", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
			mockUserClient.AssertExpectations(t)
			mockAMQPService.AssertExpectations(t)
			mockS3Service.AssertExpectations(t)
			mockMailRepo.AssertExpectations(t)
		})
	}
}
