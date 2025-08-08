package send_mail

import (
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

func TestSendMailController_GetSendMailByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		sendMailID     string
		expectedStatus int
		setupMock      func(*mocks.MockSendMailRepository, primitive.ObjectID, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:           "Success",
			sendMailID:     primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sendMailID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mail := &models.Mail{
					ID:     &mailID,
					UserID: userID,
				}
				sendMail := &models.SendMail{
					ID:         sendMailID,
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}
				mockRepo.On("GetByID", mock.Anything, sendMailID).Return(sendMail, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid send mail ID",
			sendMailID:     "invalid-id",
			expectedStatus: http.StatusBadRequest,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sendMailID primitive.ObjectID) {
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Send mail not found",
			sendMailID:     primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusNotFound,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sendMailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, sendMailID).Return(nil, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Unauthorized",
			sendMailID:     primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sendMailID primitive.ObjectID) {
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			userID := primitive.NewObjectID()

			var sendMailID primitive.ObjectID
			if tt.name != "Invalid send mail ID" {
				sendMailID, _ = primitive.ObjectIDFromHex(tt.sendMailID)
			}

			tt.setupMock(mockRepo, userID, sendMailID)

			controller := NewSendMailController(mockRepo)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			sendMailRoutes := router.Group("/mail/send")
			{
				sendMailRoutes.GET("/:id", controller.GetSendMailByID)
			}

			req, _ := http.NewRequest("GET", "/mail/send/"+tt.sendMailID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
