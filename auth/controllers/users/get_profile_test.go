package users

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetMyProfile(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockUserRoleRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successful profile retrieval",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				password := "should-be-removed"
				user := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
				}
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
				userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]models.UserEntity
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response["data"])
				assert.Equal(t, "test@example.com", *response["data"].Email)
				assert.Nil(t, response["data"].Password, "Password should be removed from response")
			},
		},
		{
			name:           "Unauthorized - no auth user",
			setupAuth:      func(c *gin.Context) {},
			setupMocks:     func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Authentication required", response["error"])
			},
		},
		{
			name: "Error fetching user profile",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to retrieve user profile", response["error"])
			},
		},
		{
			name: "Error populating roles",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				user := &models.UserEntity{ID: &userID}
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)
				userRoleRepo.On("PopulateRoles", mock.Anything, user).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to populate user roles", response["error"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockUserRepo := new(mocks.MockUserRepository)
			mockUserRoleRepo := new(mocks.MockUserRoleRepository)

			// Setup mocks
			tc.setupMocks(mockUserRepo, mockUserRoleRepo)

			// Create controller and router
			controller := NewUserController(mockUserRepo, mockUserRoleRepo, new(mocks.MockProductivityClient))
			router := gin.New()
			router.GET("/users/me", func(c *gin.Context) {
				tc.setupAuth(c)
				controller.GetMyProfile(c)
			})

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users/me", nil)

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			tc.checkResponse(t, w)

			// Verify mock expectations
			mockUserRepo.AssertExpectations(t)
			mockUserRoleRepo.AssertExpectations(t)
		})
	}
}
