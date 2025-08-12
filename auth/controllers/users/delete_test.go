package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"

	"connectrpc.com/connect"
	"github.com/atomic-blend/backend/grpc/gen/productivity/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteAccount(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockUserRoleRepository, *mocks.MockProductivityClient)
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Successful account deletion",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository, productivityClient *mocks.MockProductivityClient) {
				userID := primitive.NewObjectID()
				user := &models.UserEntity{ID: &userID}
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
				userRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

				// Mock successful productivity client call
				productivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(
					&connect.Response[productivityv1.DeleteUserDataResponse]{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Account successfully deleted"},
		},
		{
			name:      "Unauthorized - no auth user",
			setupAuth: func(c *gin.Context) {},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository, productivityClient *mocks.MockProductivityClient) {
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]string{"error": "Authentication required"},
		},
		{
			name: "User not found",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository, productivityClient *mocks.MockProductivityClient) {
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]string{"error": "User not found"},
		},
		{
			name: "Error during personal data deletion",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository, productivityClient *mocks.MockProductivityClient) {
				userID := primitive.NewObjectID()
				user := &models.UserEntity{ID: &userID}
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Mock failing productivity client call
				productivityClient.On("DeleteUserData", mock.Anything, mock.Anything).Return(
					nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "Failed to delete personal data: assert.AnError general error for testing"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockUserRepo := new(mocks.MockUserRepository)
			mockUserRoleRepo := new(mocks.MockUserRoleRepository)
			mockProductivityClient := new(mocks.MockProductivityClient)

			// Setup mocks
			tc.setupMocks(mockUserRepo, mockUserRoleRepo, mockProductivityClient)

			// Create controller and router
			controller := NewUserController(mockUserRepo, mockUserRoleRepo, mockProductivityClient)

			router := gin.New()
			router.DELETE("/users/me", func(c *gin.Context) {
				tc.setupAuth(c)
				controller.DeleteAccount(c)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/users/me", nil)

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert response body
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, response)

			// Verify mock expectations
			mockUserRepo.AssertExpectations(t)
			mockUserRoleRepo.AssertExpectations(t)
			mockProductivityClient.AssertExpectations(t)
		})
	}
}
