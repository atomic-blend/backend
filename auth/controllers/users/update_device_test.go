package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateDeviceInfo(t *testing.T) {
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
			name: "Successfully add new device",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"deviceId":   "device123",
				"deviceName": "My Test Device",
				"fcmToken":   "fcm-token-123",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				password := "password-hash"

				// User with no devices initially
				user := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices:  []*models.UserDevice{},
				}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Expected updated user with new device
				updatedUser := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices: []*models.UserDevice{
						{
							DeviceID:   "device123",
							DeviceName: "My Test Device",
							FcmToken:   "fcm-token-123",
						},
					},
				}

				// Mock updating user
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(updatedUser, nil)

				// Mock PopulateRoles call
				userRoleRepo.On("PopulateRoles", mock.Anything, updatedUser).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Device information updated successfully", response["message"])

				// Verify response data
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)

				devices, ok := data["devices"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, devices, 1)

				device := devices[0].(map[string]interface{})
				assert.Equal(t, "device123", device["deviceId"])
				assert.Equal(t, "My Test Device", device["deviceName"])
				assert.Equal(t, "fcm-token-123", device["fcmToken"])
			},
		},
		{
			name: "Successfully update existing device",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"deviceId":   "existing-device-id",
				"deviceName": "Updated Device Name",
				"fcmToken":   "updated-fcm-token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				password := "password-hash"

				// User with existing device
				user := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices: []*models.UserDevice{
						{
							DeviceID:   "existing-device-id",
							DeviceName: "Old Device Name",
							FcmToken:   "old-fcm-token",
						},
					},
				}

				// Mock finding user by ID
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Expected updated user with updated device
				updatedUser := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices: []*models.UserDevice{
						{
							DeviceID:   "existing-device-id",
							DeviceName: "Updated Device Name",
							FcmToken:   "updated-fcm-token",
						},
					},
				}

				// Mock updating user
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(updatedUser, nil)

				// Mock PopulateRoles call
				userRoleRepo.On("PopulateRoles", mock.Anything, updatedUser).Return(nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Device information updated successfully", response["message"])

				// Verify response data
				data, ok := response["data"].(map[string]interface{})
				assert.True(t, ok)

				devices, ok := data["devices"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, devices, 1)

				device := devices[0].(map[string]interface{})
				assert.Equal(t, "existing-device-id", device["deviceId"])
				assert.Equal(t, "Updated Device Name", device["deviceName"])
				assert.Equal(t, "updated-fcm-token", device["fcmToken"])
			},
		},
		{
			name:      "Unauthorized - no auth user",
			setupAuth: func(c *gin.Context) {},
			reqBody: map[string]interface{}{
				"deviceId":   "device123",
				"deviceName": "My Device",
				"fcmToken":   "fcm-token",
			},
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
			name: "Invalid request format",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				// Missing required fields
				"deviceName": "My Device",
			},
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
		{
			name: "Error fetching user",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"deviceId":   "device123",
				"deviceName": "My Device",
				"fcmToken":   "fcm-token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				// Mock error when finding user
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, errors.New("database error"))
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
			reqBody: map[string]interface{}{
				"deviceId":   "device123",
				"deviceName": "My Device",
				"fcmToken":   "fcm-token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				user := &models.UserEntity{
					ID:    &userID,
					Email: &email,
				}

				// Mock finding user successfully
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				// Mock error when updating user
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(nil, errors.New("update failed"))

				// We don't mock PopulateRoles here because Update returns an error,
				// so PopulateRoles should never be called
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Failed to update device information", response["error"])
			},
		},
		{
			name: "Error populating roles",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			reqBody: map[string]interface{}{
				"deviceId":   "device123",
				"deviceName": "My Device",
				"fcmToken":   "fcm-token",
			},
			setupMocks: func(userRepo *mocks.MockUserRepository, userRoleRepo *mocks.MockUserRoleRepository) {
				userID := primitive.NewObjectID()
				email := "test@example.com"
				password := "password-hash"
				user := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices:  []*models.UserDevice{},
				}

				// Mock finding user successfully
				userRepo.On("FindByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(user, nil)

				updatedUser := &models.UserEntity{
					ID:       &userID,
					Email:    &email,
					Password: &password,
					Devices: []*models.UserDevice{
						{
							DeviceID:   "device123",
							DeviceName: "My Device",
							FcmToken:   "fcm-token",
						},
					},
				}

				// Mock updating user successfully
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserEntity")).Return(updatedUser, nil)

				// Mock error when populating roles
				userRoleRepo.On("PopulateRoles", mock.Anything, updatedUser).Return(errors.New("role population error"))
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
			router.PUT("/users/device", func(c *gin.Context) {
				tc.setupAuth(c)
				controller.UpdateDeviceInfo(c)
			})

			// Convert request body to JSON
			bodyJSON, _ := json.Marshal(tc.reqBody)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPut, "/users/device", bytes.NewBuffer(bodyJSON))
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
