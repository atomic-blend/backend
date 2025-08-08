package send_mail

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/auth"
	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendMailController_CreateSendMail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		setupMock      func(*mocks.MockSendMailRepository, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name: "Success",
			requestBody: CreateSendMailRequest{
				Mail: &models.Mail{
					TextContent: "Test email content",
					Headers: map[string]string{
						"Subject": "Test Email",
						"From":    "test@example.com",
					},
				},
			},
			expectedStatus: http.StatusCreated,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {
				mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(sendMail *models.SendMail) bool {
					return sendMail.Mail != nil &&
						sendMail.Mail.UserID == userID &&
						sendMail.SendStatus == models.SendStatusPending &&
						sendMail.RetryCounter == nil // Should be nil when created
				})).Return(&models.SendMail{
					ID:           primitive.NewObjectID(),
					SendStatus:   models.SendStatusPending,
					RetryCounter: nil, // Should be nil when created
				}, nil)
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
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Missing mail field",
			requestBody:    CreateSendMailRequest{},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name: "Unauthorized",
			requestBody: CreateSendMailRequest{
				Mail: &models.Mail{TextContent: "Test"},
			},
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID) {},
			setupAuth:      func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			userID := primitive.NewObjectID()
			tt.setupMock(mockRepo, userID)

			controller := NewSendMailController(mockRepo)

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
		})
	}
}
