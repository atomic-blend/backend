package draftmail

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestDraftMailController_DeleteDraftMail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		draftMailID    string
		expectedStatus int
		setupMock      func(*mocks.MockDraftMailRepository, primitive.ObjectID, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:           "Success",
			draftMailID:    primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusNoContent,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: userID}
				draftMail := &models.SendMail{
					ID:         draftMailID,
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}
				mockRepo.On("GetByID", mock.Anything, draftMailID).Return(draftMail, nil)
				mockRepo.On("Delete", mock.Anything, draftMailID).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid draft mail ID",
			draftMailID:    "invalid-id",
			expectedStatus: http.StatusBadRequest,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Draft mail not found",
			draftMailID:    primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusNotFound,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, draftMailID).Return(nil, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Unauthorized",
			draftMailID:    primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID, draftMailID primitive.ObjectID) {
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {},
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

			controller := NewDraftMailController(mockRepo, mockUserClient, mockAMQPService, mockS3Service)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			draftMailRoutes := router.Group("/mail/draft")
			{
				draftMailRoutes.DELETE("/:id", controller.DeleteDraftMail)
			}

			req, _ := http.NewRequest("DELETE", "/mail/draft/"+tt.draftMailID, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
