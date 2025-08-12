package sendmail

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/mail/models"
	amqpservice "github.com/atomic-blend/backend/shared/services/amqp"
	s3service "github.com/atomic-blend/backend/shared/services/s3"
	"github.com/atomic-blend/backend/mail/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestSendMailController_DeleteSendMail(t *testing.T) {
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
			expectedStatus: http.StatusNoContent,
			setupMock: func(mockRepo *mocks.MockSendMailRepository, userID primitive.ObjectID, sendMailID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: userID}
				sendMail := &models.SendMail{
					ID:         sendMailID,
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}
				mockRepo.On("GetByID", mock.Anything, sendMailID).Return(sendMail, nil)
				mockRepo.On("Delete", mock.Anything, sendMailID).Return(nil)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockSendMailRepository{}
			mockUserClient := &mocks.MockUserClient{}
			mockAMQPService := &amqpservice.MockAMQPService{}
			mockS3Service := &s3service.MockS3Service{}
			userID := primitive.NewObjectID()

			var sendMailID primitive.ObjectID
			if tt.name != "Invalid send mail ID" {
				sendMailID, _ = primitive.ObjectIDFromHex(tt.sendMailID)
			}

			tt.setupMock(mockRepo, userID, sendMailID)

			controller := NewSendMailController(mockRepo, mockUserClient, mockAMQPService, mockS3Service)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			sendMailRoutes := router.Group("/mail/send")
			{
				sendMailRoutes.DELETE("/:id", controller.DeleteSendMail)
			}

			req, _ := http.NewRequest("DELETE", "/mail/send/"+tt.sendMailID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
