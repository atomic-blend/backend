package mail

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMailController_GetMailByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mailID         string
		expectedStatus int
		setupMock      func(*mocks.MockMailRepository, primitive.ObjectID, primitive.ObjectID)
		setupAuth      func(*gin.Context, primitive.ObjectID)
	}{
		{
			name:           "Success",
			mailID:         primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mail := &models.Mail{
					ID:     &mailID,
					UserID: userID,
					Headers: map[string]string{
						"Subject": "Test Email",
						"From":    "test@example.com",
					},
				}
				mockRepo.On("GetByID", mock.Anything, mailID).Return(mail, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Invalid mail ID",
			mailID:         "invalid-id",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
		{
			name:           "Mail not found",
			mailID:         primitive.NewObjectID().Hex(),
			expectedStatus: http.StatusNotFound,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID, mailID primitive.ObjectID) {
				mockRepo.On("GetByID", mock.Anything, mailID).Return(nil, nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mocks.MockMailRepository{}
			userID := primitive.NewObjectID()

			// Parse the mailID from the test case for consistent usage
			var mailID primitive.ObjectID
			if tt.name != "Invalid mail ID" {
				mailID, _ = primitive.ObjectIDFromHex(tt.mailID)
			}

			tt.setupMock(mockRepo, userID, mailID)

			controller := NewMailController(mockRepo)

			// Create router and add auth middleware
			router := gin.New()
			router.Use(func(c *gin.Context) {
				tt.setupAuth(c, userID)
				c.Next()
			})

			// Setup routes
			mailRoutes := router.Group("/mail")
			{
				mailRoutes.GET("/:id", controller.GetMailByID)
			}

			// Create request
			req, _ := http.NewRequest("GET", "/mail/"+tt.mailID, nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			mockRepo.AssertExpectations(t)
		})
	}
}
