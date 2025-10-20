package waitinglist

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJoinWaitingList(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMocks     func(*mocks.MockWaitingListRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful join waiting list",
			requestBody: JoinWaitingListRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

				// Mock Count to return 5 (5 people already in waiting list)
				mockRepo.On("Count", mock.Anything).Return(int64(5), nil)

				// Mock Create to return a new waiting list record
				now := primitive.NewDateTimeFromTime(time.Now())
				expectedRecord := &waitinglist.WaitingList{
					Email:     "test@example.com",
					CreatedAt: &now,
					UpdatedAt: &now,
				}
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(expectedRecord, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":       "success",
				"before_count": float64(5),
			},
		},
		{
			name: "email already in waiting list",
			requestBody: JoinWaitingListRequest{
				Email: "existing@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
				// Mock GetByEmail to return existing record
				now := primitive.NewDateTimeFromTime(time.Now())
				existingRecord := &waitinglist.WaitingList{
					Email:     "existing@example.com",
					CreatedAt: &now,
					UpdatedAt: &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingRecord, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "error_email_already_in_waiting_list",
			},
		},
		{
			name: "invalid email format",
			requestBody: JoinWaitingListRequest{
				Email: "invalid-email",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
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
				"invalid_field": "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
				// No mocks needed as validation happens before repository calls
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": mock.AnythingOfType("string"),
			},
		},
		{
			name: "database error when checking existing email",
			requestBody: JoinWaitingListRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
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
			requestBody: JoinWaitingListRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

				// Mock Count to return error
				mockRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "error_getting_waiting_list_count",
			},
		},
		{
			name: "database error when creating record",
			requestBody: JoinWaitingListRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

				// Mock Count to return 0
				mockRepo.On("Count", mock.Anything).Return(int64(0), nil)

				// Mock Create to return error
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "error_creating_waiting_list_record",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mocks
			mockRepo := new(mocks.MockWaitingListRepository)

			// Setup mocks
			tc.setupMocks(mockRepo)

			// Create controller and router
			controller := NewController(mockRepo)

			router := gin.New()
			router.POST("/waiting-list", controller.JoinWaitingList)

			// Create request body
			requestBodyBytes, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/waiting-list", bytes.NewBuffer(requestBodyBytes))
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
				assert.Equal(t, tc.expectedBody["before_count"], response["before_count"])
				assert.Contains(t, response, "entry")
			} else {
				assert.Contains(t, response, "error")
				if tc.expectedBody["error"] != mock.AnythingOfType("string") {
					assert.Equal(t, tc.expectedBody["error"], response["error"])
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestJoinWaitingListIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should work with real database", func(t *testing.T) {
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

		// Create controller with real repository
		// Note: controller will be created by SetupRoutes with real repository

		// Create router
		router := gin.New()
		SetupRoutes(router, db)

		// Test successful join
		requestBody := JoinWaitingListRequest{
			Email: "integration@example.com",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/waiting-list", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// Should return 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response["status"])
		assert.Contains(t, response, "entry")
		assert.Contains(t, response, "before_count")

		// Test duplicate email
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPost, "/waiting-list", bytes.NewBuffer(requestBodyBytes))
		req2.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w2, req2)

		// Should return 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		// Parse error response
		var errorResponse map[string]string
		err = json.Unmarshal(w2.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error_email_already_in_waiting_list", errorResponse["error"])
	})
}

func TestJoinWaitingListRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	validationTestCases := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name: "valid email",
			requestBody: JoinWaitingListRequest{
				Email: "valid@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "email with subdomain",
			requestBody: JoinWaitingListRequest{
				Email: "user@subdomain.example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "email with plus sign",
			requestBody: JoinWaitingListRequest{
				Email: "user+tag@example.com",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid email - no @",
			requestBody: JoinWaitingListRequest{
				Email: "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email - no domain",
			requestBody: JoinWaitingListRequest{
				Email: "user@",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid email - no local part",
			requestBody: JoinWaitingListRequest{
				Email: "@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty email",
			requestBody: JoinWaitingListRequest{
				Email: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "email with spaces",
			requestBody: JoinWaitingListRequest{
				Email: "user @example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range validationTestCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := new(mocks.MockWaitingListRepository)

			// Setup mocks for successful cases
			if tc.expectedStatus == http.StatusOK {
				mockRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
				mockRepo.On("Count", mock.Anything).Return(int64(0), nil)
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(&waitinglist.WaitingList{}, nil)
			}

			// Create controller and router
			controller := NewController(mockRepo)

			router := gin.New()
			router.POST("/waiting-list", controller.JoinWaitingList)

			// Create request body
			requestBodyBytes, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/waiting-list", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}
