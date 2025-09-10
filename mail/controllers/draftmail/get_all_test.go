package draftmail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	amqpservice "github.com/atomic-blend/backend/shared/services/amqp"
	s3service "github.com/atomic-blend/backend/shared/services/s3"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDraftMailController_GetAllDraftMails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		setupMock      func(*mocks.MockDraftMailRepository, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:           "Success with default pagination",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mail := &models.Mail{ID: &mailID, UserID: userID}
				draftMails := []*models.SendMail{{
					ID:         primitive.NewObjectID(),
					Mail:       mail,
					SendStatus: models.SendStatusPending,
				}}
				mockRepo.On("GetAll", mock.Anything, userID, int64(1), int64(10)).Return(draftMails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Success with custom pagination",
			queryParams:    "?page=2&size=15",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID) {
				draftMails := []*models.SendMail{}
				mockRepo.On("GetAll", mock.Anything, userID, int64(2), int64(15)).Return(draftMails, int64(0), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Unauthorized",
			queryParams:    "",
			expectedStatus: http.StatusUnauthorized,
			setupMock:      func(mockRepo *mocks.MockDraftMailRepository, userID primitive.ObjectID) {},
			setupAuth:      func(c *gin.Context, userID primitive.ObjectID) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockDraftMailRepository{}
			mockUserClient := &mocks.MockUserClient{}
			mockAMQPService := &amqpservice.MockAMQPService{}
			mockS3Service := &s3service.MockS3Service{}
			userID := primitive.NewObjectID()
			tt.setupMock(mockRepo, userID)

			controller := NewDraftMailController(mockRepo, mockUserClient, mockAMQPService, mockS3Service)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			draftMailRoutes := router.Group("/mail/draft")
			{
				draftMailRoutes.GET("", pagination.New(), controller.GetAllDraftMails)
			}

			req, _ := http.NewRequest("GET", "/mail/draft"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response PaginatedDraftMailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.DraftMails)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
