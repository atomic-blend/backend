package mail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/webstradev/gin-pagination/v2/pkg/pagination"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMailController_GetAllMails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		setupMock      func(*mocks.MockMailRepository, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:           "Success with default pagination",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mails := []*models.Mail{{ID: &mailID, UserID: userID}}
				mockRepo.On("GetAll", mock.Anything, userID, int64(1), int64(10)).Return(mails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Success with custom pagination",
			queryParams:    "?page=2&size=15",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mails := []*models.Mail{}
				mockRepo.On("GetAll", mock.Anything, userID, int64(2), int64(15)).Return(mails, int64(0), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Success with page only",
			queryParams:    "?page=3",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mails := []*models.Mail{{ID: &mailID, UserID: userID}}
				mockRepo.On("GetAll", mock.Anything, userID, int64(3), int64(10)).Return(mails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Success with size only",
			queryParams:    "?size=15",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mailID := primitive.NewObjectID()
				mails := []*models.Mail{{ID: &mailID, UserID: userID}}
				mockRepo.On("GetAll", mock.Anything, userID, int64(1), int64(15)).Return(mails, int64(1), nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.MockMailRepository{}
			userID := primitive.NewObjectID()
			tt.setupMock(mockRepo, userID)

			controller := NewMailController(mockRepo)

			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			mailRoutes := router.Group("/mail")
			{
				mailRoutes.GET("", pagination.New(), controller.GetAllMails)
			}

			req, _ := http.NewRequest("GET", "/mail"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response PaginatedMailResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response.Mails)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
