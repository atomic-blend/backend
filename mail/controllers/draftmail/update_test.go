package draftmail

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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDraftMailController_UpdateDraftMail(t *testing.T) {
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
		draftMailID    string
		requestBody    interface{}
		expectedStatus int
		setupMock      func(*mocks.MockDraftMailRepository, primitive.ObjectID, primitive.ObjectID)
		setupUserMock  func(*mocks.MockUserClient, primitive.ObjectID)
		setupS3Mock    func(*s3service.MockS3Service, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:        "Success",
			draftMailID: primitive.NewObjectID().Hex(),
			requestBody: models.RawMail{
				Headers: map[string]interface{}{
					"Subject": "Updated Test Draft Email",
					"From":    "test@example.com",
				},
				TextContent: "Updated test draft email content",
			},
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: userID}
				existingDraftMail := &models.SendMail{
					ID:         draftMailID,
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}
				updatedDraftMail := &models.SendMail{
					ID:         draftMailID,
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}
				mockRepo.On("GetByID", mock.Anything, draftMailID).Return(existingDraftMail, nil)
				mockRepo.On("Update", mock.Anything, draftMailID, mock.Anything).Return(updatedDraftMail, nil)
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
			setupS3Mock: func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {
				// No S3 operations expected since there are no attachments in this test
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid draft mail ID",
			draftMailID:    "invalid-id",
			requestBody:    models.RawMail{TextContent: "Test"},
			expectedStatus: http.StatusBadRequest,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
			},
			setupUserMock: func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupS3Mock:   func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Draft mail not found",
			draftMailID:    primitive.NewObjectID().Hex(),
			requestBody:    models.RawMail{TextContent: "Test"},
			expectedStatus: http.StatusNotFound,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, draftMailID).Return(nil, nil)
			},
			setupUserMock: func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupS3Mock:   func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Unauthorized",
			draftMailID:    primitive.NewObjectID().Hex(),
			requestBody:    models.RawMail{TextContent: "Test"},
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
			},
			setupUserMock: func(mockUserClient *mocks.MockUserClient, userID primitive.ObjectID) {},
			setupS3Mock:   func(mockS3Service *s3service.MockS3Service, userID primitive.ObjectID) {},
			setupAuth:     func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockDraftMailRepository{}
			mockUserClient := &mocks.MockUserClient{}
			mockAMQPService := &amqpservice.MockAMQPService{}
			mockS3Service := &s3service.MockS3Service{}
			userID := primitive.NewObjectID()

			var draftMailID primitive.ObjectID
			if tt.name != "Invalid draft mail ID" {
				draftMailID, _ = primitive.ObjectIDFromHex(tt.draftMailID)
			}

			tt.setupMock(mockRepo, userID, draftMailID)
			tt.setupUserMock(mockUserClient, userID)
			tt.setupS3Mock(mockS3Service, userID)

			controller := NewDraftMailController(mockRepo, mockUserClient, mockAMQPService, mockS3Service)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			draftMailRoutes := router.Group("/mail/draft")
			{
				draftMailRoutes.PUT("/:id", controller.UpdateDraftMail)
			}

			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req, _ := http.NewRequest("PUT", "/mail/draft/"+tt.draftMailID, &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
			mockUserClient.AssertExpectations(t)
			mockS3Service.AssertExpectations(t)
		})
	}
}
