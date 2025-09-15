package mail

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/mail/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMailController_CleanupTrash(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		expectedStatus   int
		setupMock        func(*mocks.MockMailRepository, primitive.ObjectID)
		setupAuth        func(*gin.Context, primitive.ObjectID)
		expectedResponse map[string]interface{}
	}{
		{
			name:           "Success cleanup trash",
			expectedStatus: http.StatusOK,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mockRepo.On("CleanupTrash", mock.Anything, &userID).Return(nil)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			expectedResponse: map[string]interface{}{
				"message": "Trash cleanup completed successfully",
			},
		},
		{
			name:           "Repository error during cleanup",
			expectedStatus: http.StatusInternalServerError,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				mockRepo.On("CleanupTrash", mock.Anything, &userID).Return(assert.AnError)
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			expectedResponse: map[string]interface{}{
				"error": "Failed to cleanup trash",
			},
		},
		{
			name:           "Unauthorized - no auth user",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				// No mock expectations since auth should fail first
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				// Don't set auth user to simulate unauthorized request
			},
			expectedResponse: map[string]interface{}{
				"error": "Authentication required",
			},
		},
		{
			name:           "Unauthorized - nil auth user",
			expectedStatus: http.StatusUnauthorized,
			setupMock: func(mockRepo *mocks.MockMailRepository, userID primitive.ObjectID) {
				// No mock expectations since auth should fail first
			},
			setupAuth: func(c *gin.Context, userID primitive.ObjectID) {
				c.Set("authUser", nil)
			},
			expectedResponse: map[string]interface{}{
				"error": "Authentication required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := &mocks.MockMailRepository{}
			userID := primitive.NewObjectID()
			tt.setupMock(mockRepo, userID)

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
				mailRoutes.POST("/trash/empty", controller.CleanupTrash)
			}

			// Create request
			req, _ := http.NewRequest("POST", "/mail/trash/empty", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Check response body
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedResponse, response)

			mockRepo.AssertExpectations(t)
		})
	}
}
