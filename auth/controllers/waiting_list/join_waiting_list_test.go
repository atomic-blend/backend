package waitinglist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	amqpservice "github.com/atomic-blend/backend/shared/services/amqp"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestJoinWaitingList(t *testing.T) {
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
		setupMocks     func(*mocks.MockWaitingListRepository, *amqpservice.MockAMQPService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful join waiting list",
			requestBody: JoinWaitingListRequest{
				Email: "test@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

				// Mock Create to return a new waiting list record
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(expectedRecord, nil)

				// Mock GetPositionByEmail to return position
				mockRepo.On("GetPositionByEmail", mock.Anything, "test@example.com").Return(int64(5), nil)

				// Mock Count to return total count after creation
				mockRepo.On("Count", mock.Anything).Return(int64(6), nil)

				// Mock AMQP message publishing
				mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Return()
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"status":   "success",
				"position": float64(5),
				"total":    float64(6),
			},
		},
		{
			name: "email already in waiting list",
			requestBody: JoinWaitingListRequest{
				Email: "existing@example.com",
			},
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
				// Mock GetByEmail to return existing record
				now := primitive.NewDateTimeFromTime(time.Now())
				existingRecord := &waitinglist.WaitingList{
					Email:     "existing@example.com",
					CreatedAt: &now,
					UpdatedAt: &now,
				}
				mockRepo.On("GetByEmail", mock.Anything, "existing@example.com").Return(existingRecord, nil)
				// No AMQP message publishing expected since user already exists
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

				// Mock Create to return a new waiting list record
				now := primitive.NewDateTimeFromTime(time.Now())
				securityToken := "a1b2c3d4e5f678901234567890123456"
				expectedRecord := &waitinglist.WaitingList{
					Email:         "test@example.com",
					SecurityToken: &securityToken,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(expectedRecord, nil)

				// Mock GetPositionByEmail to return position
				mockRepo.On("GetPositionByEmail", mock.Anything, "test@example.com").Return(int64(0), nil)

				// Mock Count to return error
				mockRepo.On("Count", mock.Anything).Return(int64(0), assert.AnError)

				// Mock AMQP message publishing
				mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Return()
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
			setupMocks: func(mockRepo *mocks.MockWaitingListRepository, mockAMQP *amqpservice.MockAMQPService) {
				// Mock GetByEmail to return nil (email not found)
				mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)

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
			mockAMQP := new(amqpservice.MockAMQPService)
			// Setup mocks
			tc.setupMocks(mockRepo, mockAMQP)

			// Create controller and router
			controller := NewController(mockRepo, mockAMQP)

			router := gin.New()
			router.POST("/auth/waiting-list", controller.JoinWaitingList)

			// Create request body
			requestBodyBytes, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
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

				// Check that entry contains security token
				entry, ok := response["entry"].(map[string]interface{})
				assert.True(t, ok, "entry should be a map")
				assert.Contains(t, entry, "securityToken")
				assert.NotNil(t, entry["securityToken"])

				// Verify security token is 32 characters long
				securityToken, ok := entry["securityToken"].(string)
				assert.True(t, ok, "securityToken should be a string")
				assert.Equal(t, 32, len(securityToken), "securityToken should be 32 characters long")
			} else {
				assert.Contains(t, response, "error")
				if tc.expectedBody["error"] != mock.AnythingOfType("string") {
					assert.Equal(t, tc.expectedBody["error"], response["error"])
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockAMQP.AssertExpectations(t)
		})
	}
}

func TestJoinWaitingListIntegration(t *testing.T) {
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

	t.Run("should work with real database and email templates", func(t *testing.T) {
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

		// Create mock AMQP service
		mockAMQP := &amqpservice.MockAMQPService{}

		// Mock AMQP message publishing
		mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Return()

		// Create controller with mock AMQP service
		controller := NewController(waitingListRepo, mockAMQP)

		// Create router
		router := gin.New()
		router.POST("/auth/waiting-list", controller.JoinWaitingList)

		// Test successful join
		requestBody := JoinWaitingListRequest{
			Email: "integration@example.com",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
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
		assert.Contains(t, response, "position")
		assert.Contains(t, response, "total")

		// Verify security token in entry
		entry, ok := response["entry"].(map[string]interface{})
		assert.True(t, ok, "entry should be a map")
		assert.Contains(t, entry, "securityToken")

		securityToken, ok := entry["securityToken"].(string)
		assert.True(t, ok, "securityToken should be a string")
		assert.Equal(t, 32, len(securityToken), "securityToken should be 32 characters long")

		// Test duplicate email
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
		req2.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w2, req2)

		// Should return 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		// Parse error response
		var errorResponse map[string]string
		err = json.Unmarshal(w2.Body.Bytes(), &errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, "error_email_already_in_waiting_list", errorResponse["error"])

		// Verify mock expectations
		mockAMQP.AssertExpectations(t)
	})
}

func TestEmailTemplateRendering(t *testing.T) {
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

	t.Run("should render email templates with correct data", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(mocks.MockWaitingListRepository)
		mockAMQP := new(amqpservice.MockAMQPService)

		// Setup mocks
		mockRepo.On("GetByEmail", mock.Anything, "template@example.com").Return(nil, nil)
		mockRepo.On("Count", mock.Anything).Return(int64(0), nil)

		now := primitive.NewDateTimeFromTime(time.Now())
		securityToken := "test123456789012345678901234567890"
		expectedRecord := &waitinglist.WaitingList{
			Email:         "template@example.com",
			SecurityToken: &securityToken,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		}
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(expectedRecord, nil)

		// Mock GetPositionByEmail to return position
		mockRepo.On("GetPositionByEmail", mock.Anything, "template@example.com").Return(int64(0), nil)

		// Mock Count to return total count after creation
		mockRepo.On("Count", mock.Anything).Return(int64(1), nil)

		// Mock AMQP message publishing and capture the message to verify content
		var capturedMessage map[string]interface{}
		mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Run(func(args mock.Arguments) {
			capturedMessage = args.Get(2).(map[string]interface{})
		}).Return()

		// Create controller and router
		controller := NewController(mockRepo, mockAMQP)

		router := gin.New()
		router.POST("/auth/waiting-list", controller.JoinWaitingList)

		// Create request body
		requestBody := JoinWaitingListRequest{
			Email: "template@example.com",
		}
		requestBodyBytes, err := json.Marshal(requestBody)
		assert.NoError(t, err)

		// Create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Should return 200 OK
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify AMQP message was published with correct content
		mockAMQP.AssertExpectations(t)
		assert.NotNil(t, capturedMessage, "AMQP message should have been captured")

		// Verify AMQP message structure
		assert.Equal(t, true, capturedMessage["waiting_list_email"])
		content := capturedMessage["content"].(map[string]interface{})

		// Verify headers
		headers := content["headers"].(map[string]interface{})
		assert.Equal(t, []string{"template@example.com"}, headers["To"])
		assert.Equal(t, "noreply@atomic-blend.com", headers["From"])
		assert.Equal(t, "You just joined the waiting list!", headers["Subject"])

		// Extract the actual security token from the HTML content
		// The token is generated randomly, so we need to extract it from the response
		var actualSecurityToken string
		htmlContent := content["htmlContent"].(string)
		if len(htmlContent) > 0 {
			// Find the security token in the HTML content
			// Look for the pattern: <p style="...">TOKEN</p>
			start := strings.Index(htmlContent, "<p style=\"margin: 0; word-break: break-all; font-family: ui-monospace")
			if start != -1 {
				// Find the closing tag
				end := strings.Index(htmlContent[start:], "</p>")
				if end != -1 {
					tokenStart := strings.Index(htmlContent[start:start+end], ">")
					if tokenStart != -1 {
						tokenEnd := start + end
						tokenStartPos := start + tokenStart + 1
						actualSecurityToken = strings.TrimSpace(htmlContent[tokenStartPos:tokenEnd])
					}
				}
			}
		}

		// Verify HTML content contains email and security token
		assert.Contains(t, htmlContent, "template@example.com")
		assert.Contains(t, htmlContent, actualSecurityToken)
		assert.Contains(t, htmlContent, "Welcome to Atomic Blend!")

		// Verify text content contains email and security token
		textContent := content["textContent"].(string)
		assert.Contains(t, textContent, "template@example.com")
		assert.Contains(t, textContent, actualSecurityToken)
		assert.Contains(t, textContent, "Welcome to Atomic Blend!")

		// Verify the security token is 32 characters long
		assert.Equal(t, 32, len(actualSecurityToken), "securityToken should be 32 characters long")

		// Verify mock expectations
		mockRepo.AssertExpectations(t)
	})
}

func TestJoinWaitingListRequestValidation(t *testing.T) {
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
			mockAMQP := new(amqpservice.MockAMQPService)
			// Setup mocks for successful cases
			if tc.expectedStatus == http.StatusOK {
				mockRepo.On("GetByEmail", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)
				mockRepo.On("Count", mock.Anything).Return(int64(0), nil)
				securityToken := "a1b2c3d4e5f678901234567890123456"
				mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*waitinglist.WaitingList")).Return(&waitinglist.WaitingList{
					SecurityToken: &securityToken,
				}, nil)

				// Mock GetPositionByEmail to return position
				mockRepo.On("GetPositionByEmail", mock.Anything, mock.AnythingOfType("string")).Return(int64(0), nil)

				// Mock Count to return total count after creation
				mockRepo.On("Count", mock.Anything).Return(int64(1), nil)

				// Mock AMQP message publishing for successful cases
				mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Return()
			}

			// Create controller and router
			controller := NewController(mockRepo, mockAMQP)

			router := gin.New()
			router.POST("/auth/waiting-list", controller.JoinWaitingList)

			// Create request body
			requestBodyBytes, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockAMQP.AssertExpectations(t)
		})
	}
}

func TestSecurityTokenGeneration(t *testing.T) {
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

	t.Run("should generate unique 32-character security tokens", func(t *testing.T) {
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

		// Create mock AMQP service
		mockAMQP := &amqpservice.MockAMQPService{}

		// Mock AMQP message publishing for all requests
		mockAMQP.On("PublishMessage", "mail", "sent", mock.AnythingOfType("map[string]interface {}"), mock.AnythingOfType("*amqp.Table")).Return()

		// Create controller with mock AMQP service
		controller := NewController(waitingListRepo, mockAMQP)

		// Create router
		router := gin.New()
		router.POST("/auth/waiting-list", controller.JoinWaitingList)

		// Generate multiple requests to test token uniqueness
		tokens := make(map[string]bool)

		for i := 0; i < 10; i++ {
			requestBody := JoinWaitingListRequest{
				Email: fmt.Sprintf("test%d@example.com", i),
			}
			requestBodyBytes, err := json.Marshal(requestBody)
			assert.NoError(t, err)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/auth/waiting-list", bytes.NewBuffer(requestBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			// Should return 200 OK
			assert.Equal(t, http.StatusOK, w.Code)

			// Parse response
			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Extract security token
			entry, ok := response["entry"].(map[string]interface{})
			assert.True(t, ok, "entry should be a map")

			securityToken, ok := entry["securityToken"].(string)
			assert.True(t, ok, "securityToken should be a string")
			assert.Equal(t, 32, len(securityToken), "securityToken should be 32 characters long")

			// Check for uniqueness
			assert.False(t, tokens[securityToken], "security tokens should be unique")
			tokens[securityToken] = true
		}

		// Verify we got 10 unique tokens
		assert.Equal(t, 10, len(tokens), "should generate 10 unique tokens")

		// Verify mock expectations
		mockAMQP.AssertExpectations(t)
	})
}
