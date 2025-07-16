package users

import (
	"github.com/atomic-blend/backend/auth/auth"
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateProfile(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		reqBody        map[string]interface{}
		setupMocks     func(*mocks.MockUserRepository, *mocks.MockUserRoleRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successfully update email",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"email": "newemail@example.com",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				oldEmail := "old@example.com"
				password := "password-hash"

				user := &models.UserEntity{
					ID:       &userID,
					Email:    &oldEmail,
					Password: &password,
				}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Mock checking for existing email
				userRepo.On("FindByEmail", mock.Anything, "newemail@example.com").Return(nil, errors.New("not found"))

				// Mock updating user
				updatedUser := &models.UserEntity{
					ID:       &userID,
					Password: &password,
				}
				newEmail := "newemail@example.com"
				updatedUser.Email = &newEmail

				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(updatedUser, nil)

				// Mock populating roles
				userRoleRepo.On("PopulateRoles", mock.Anything, updatedUser).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Profile updated successfully", response["message"])
				assert.NotNil(t, response["data"])
			},
		},
		{
			name: "Error - Email already in use",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"email": "existing@example.com",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				oldEmail := "old@example.com"
				password := "password-hash"

				user := &models.UserEntity{
					ID:       &userID,
					Email:    &oldEmail,
					Password: &password,
				}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Mock finding existing email that belongs to another user
				otherUserID := primitive.NewObjectID()
				otherEmail := "existing@example.com"
				otherUser := &models.UserEntity{
					ID:    &otherUserID,
					Email: &otherEmail,
				}

				userRepo.On("FindByEmail", mock.Anything, "existing@example.com").Return(otherUser, nil)

				// No need to mock PopulateRoles since the update is never performed due to conflict
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Email is already in use", response["error"])
			},
		},
		{
			name:      "Unauthorized - no auth user",
			setupAuth: func(c *gin.Context) {},
			reqBody:   map[string]interface{}{"email": "test@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				// No mock calls expected for userRepo
				// No mock calls expected for userRoleRepo
			},
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
			reqBody: map[string]interface{}{"email": "test@example.com"},
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
			name: "Error updating user",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{"email": "new@example.com"},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				oldEmail := "old@example.com"
				user := &models.UserEntity{ID: &userID, Email: &oldEmail}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Mock finding no existing email
				userRepo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, errors.New("not found"))

				// Mock update failure
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(nil, errors.New("update failed"))

				// No need to mock PopulateRoles since Update fails
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to update profile", response["error"])
			},
		},
		{
			name: "Invalid request format",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{"email": "not-an-email"},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				// No mock calls expected
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Invalid request format", response["error"])
				assert.NotNil(t, response["details"])
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
			controller := NewUserController(mockUserRepo, mockUserRoleRepo)
			router := gin.New()
			router.PUT("/users/profile", func(c *gin.Context) {
				tc.setupAuth(c)
				controller.UpdateProfile(c)
			})

			// Convert request body to JSON
			bodyJSON, _ := json.Marshal(tc.reqBody)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/users/profile", bytes.NewBuffer(bodyJSON))
			req.Header.Set("Content-Type", "application/json")

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
