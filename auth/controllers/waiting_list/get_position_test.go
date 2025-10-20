package waitinglist

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetWaitingListPosition(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Change to the auth directory to find email templates
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to auth directory where templates are located
	err = os.Chdir("../..")
	if err != nil {
		t.Fatalf("Failed to change to auth directory: %v", err)
	}

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockWaitingListRepository, *mocks.MockMailServerClient)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful get position",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return existing record
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedRecord, nil)

				// Mock Count to return total count
				mockRepo.On("Count", mock.Anything).Return(int64(10), nil)

				// Mock GetPositionByEmail to return position
				mockRepo.On("GetPositionByEmail", mock.Anything, "test@example.com").Return(int64(5), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":   "success",
				"position": float64(5),
				"total":    float64(10),
			},
		},
		{
			name: "email not found in waiting list",
			requestBody: GetWaitingListPositionRequest{
				Email:         "notfound@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "error_email_not_found_in_waiting_list",
			},
		},
		{
			name: "invalid security token",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "invalid_token",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return existing record with different token
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedRecord, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "error_invalid_security_token",
			},
		},
		{
			name: "missing security token in record",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return record without security token
				now := primitive.NewDateTimeFromTime(time.Now())
				expectedRecord := &waitinglist.WaitingList{
					Email:     "test@example.com",
					CreatedAt: &now,
					UpdatedAt: &now,
					// SecurityToken is nil
				}
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedRecord, nil)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "error_invalid_security_token",
			},
		},
		{
			name: "invalid email format",
			requestBody: GetWaitingListPositionRequest{
				Email:         "invalid-email",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// No mocks needed as validation happens before repository calls
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name: "missing email field",
			requestBody: map[string]string{
				"securityToken": "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// No mocks needed as validation happens before repository calls
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name: "missing security token field",
			requestBody: map[string]string{
				"email": "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// No mocks needed as validation happens before repository calls
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name:        "empty request body",
			requestBody: map[string]string{},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// No mocks needed as validation happens before repository calls
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name: "database error when getting record",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return error
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "error_getting_waiting_list_record",
			},
		},
		{
			name: "database error when getting count",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return existing record
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedRecord, nil)

				// Mock Count to return error
				mockRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "error_getting_waiting_list_count",
			},
		},
		{
			name: "database error when getting position",
			requestBody: GetWaitingListPositionRequest{
				Email:         "test@example.com",
				SecurityToken: "a1b2c3d4e5f678901234567890123456",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockMailClient *mocks.MockMailServerClient) {
				// Mock GetByEmail to return existing record
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(expectedRecord, nil)

				// Mock Count to return total count
				mockRepo.On("Count", mock.Anything).Return(int64(10), nil)

				// Mock GetPositionByEmail to return error
				mockRepo.On("GetPositionByEmail", mock.Anything, "test@example.com").Return(int64(0), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "error_getting_waiting_list_position",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(mocks.MockWaitingListRepository)
			mockMailServerClient := new(mocks.MockMailServerClient)
			// Setup mocks
			tc.setupMocks(mockRepo, mockMailServerClient)

			// Create controller and router
			controller := NewController(mockRepo, mockMailServerClient)

			router := gin.New()
			router.POST("/auth/waiting-list/position", controller.GetWaitingListPosition)

			// Create request body
			requestBodyBytes, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list/position", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert response body
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check specific fields based on expected status
			if tc.expectedStatus == http.StatusOK {
				assert.Equal(t, tc.expectedBody["status"], response["status"])
				assert.Equal(t, tc.expectedBody["position"], response["position"])
				assert.Equal(t, tc.expectedBody["total"], response["total"])
				assert.Contains(t, response, "entry")

				// Check that entry contains the expected fields
				entry, ok := response["entry"].(map[string]interface{})
				assert.True(t, ok, "entry should be a map")
				assert.Contains(t, entry, "email")
				assert.Contains(t, entry, "securityToken")
			} else {
				assert.Contains(t, response, "error")
				if tc.expectedBody["error"] != mock.AnythingOfType("string") {
					assert.Equal(t, tc.expectedBody["error"], response["error"])
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockMailServerClient.AssertExpectations(t)
		})
	}
}

func TestGetWaitingListPositionIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Change to the auth directory to find email templates
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(originalDir)

	// Change to auth directory where templates are located
	err = os.Chdir("../..")
	if err != nil {
		t.Fatalf("Failed to change to auth directory: %v", err)
	}

	t.Run("should work with real database and return correct position", func(t *testing.T) {
		// Start in-memory MongoDB server
		mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
		if err != nil {
			t.Fatalf("Failed to create in-memory MongoDB: %v", err)
		}
		defer mongoServer.Stop()

		// Connect to the in-memory MongoDB
		client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
		if err != nil {
			t.Fatalf("Failed to connect to in-memory MongoDB: %v", err)
		}
		defer client.Disconnect(context.Background())

		// Get database reference
		db := client.Database("test_db")

		// Create repositories
		waitingListRepo := repositories.NewWaitingListRepository(db)

		// Create mock mail server client
		mockMailServerClient := &mocks.MockMailServerClient{}

		// Create controller with mock mail server
		controller := NewController(waitingListRepo, mockMailServerClient)

		// Create router
		router := gin.New()
		router.POST("/auth/waiting-list/position", controller.GetWaitingListPosition)

		// First, create some waiting list entries to test position calculation
		emails := []string{"first@example.com", "second@example.com", "third@example.com"}
		securityTokens := []string{"token1", "token2", "token3"}

		for i, email := range emails {
			// Create waiting list entry
			now := primitive.NewDateTimeFromTime(time.Now().Add(time.Duration(i) * time.Second))
			securityToken := securityTokens[i]
			waitingListRecord := &waitinglist.WaitingList{
				Email:         email,
				SecurityToken: &securityToken,
				CreatedAt:     &now,
				UpdatedAt:     &now,
			}
			_, err := waitingListRepo.Create(context.Background(), waitingListRecord)
			assert.NoError(t, err)
		}

		// Test getting position for the second entry (should be position 1)
		requestBody := GetWaitingListPositionRequest{
			Email:         "second@example.com",
			SecurityToken: "token2",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list/position", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// Should return 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.Equal(t, float64(1), response["position"]) // Second entry should be at position 1
		assert.Equal(t, float64(3), response["total"])    // Total should be 3
		assert.Contains(t, response, "entry")

		// Test getting position for the first entry (should be position 0)
		requestBody2 := GetWaitingListPositionRequest{
			Email:         "first@example.com",
			SecurityToken: "token1",
		}
		requestBodyBytes2, err := json.Marshal(requestBody2)
		assert.NoError(t, err)

		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list/position", bytes.NewBuffer(requestBodyBytes2))
		req2.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w2, req2)

		// Should return 200 OK
		assert.Equal(t, http.StatusOK, w2.Code)

		// Parse response
		var response2 map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &response2)
		assert.NoError(t, err)
		assert.Equal(t, "success", response2["status"])
		assert.Equal(t, float64(0), response2["position"]) // First entry should be at position 0
		assert.Equal(t, float64(3), response2["total"])    // Total should be 3

		// Test with invalid security token
		requestBody3 := GetWaitingListPositionRequest{
			Email:         "first@example.com",
			SecurityToken: "invalid_token",
		}
		requestBodyBytes3, err := json.Marshal(requestBody3)
		assert.NoError(t, err)

		w3 := httptest.NewRecorder()
		req3, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list/position", bytes.NewBuffer(requestBodyBytes3))
		req3.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w3, req3)

		// Should return 401 Unauthorized
		assert.Equal(t, http.StatusUnauthorized, w3.Code)

		// Parse error response
		var errorResponse map[string]string
		err = json.Unmarshal(w3.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error_invalid_security_token", errorResponse["error"])

		// Verify mock expectations
		mockMailServerClient.AssertExpectations(t)
	})
}
