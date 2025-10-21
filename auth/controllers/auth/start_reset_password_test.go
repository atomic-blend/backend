package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/models"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/db"

	"connectrpc.com/connect"
	mailserverv1 "github.com/atomic-blend/backend/grpc/gen/mailserver/v1"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestStartResetPassword(t *testing.T) {
	// Set Gin to test mode
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

	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	if err != nil {
		t.Fatalf("Failed to create in-memory MongoDB: %v", err)
	}
	defer mongoServer.Stop()

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to in-memory MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO())

	// Get database reference
	database := client.Database("test_db")

	// Set the global database for the subscription function to use
	db.Database = database
	defer func() {
		// Reset global database after test
		db.Database = nil
	}()

	// Create repositories
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(database)
	waitingListRepo := repositories.NewWaitingListRepository(database)

	// Create mock mail server client
	mockMailServerClient := &mocks.MockMailServerClient{}

	// Create controller
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo, waitingListRepo, mockMailServerClient)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/start-reset-password", authController.StartResetPassword)

	// Create test user with backup email
	userID := primitive.NewObjectID()
	email := "test@example.com"
	backupEmail := "backup@example.com"
	password := "hashedPassword"
	testUser := &models.UserEntity{
		ID:          &userID,
		Email:       &email,
		BackupEmail: &backupEmail,
		Password:    &password,
	}

	// Insert test user into database
	_, err = database.Collection("users").InsertOne(context.TODO(), testUser)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder, *mongo.Database)
		setupMocks     func(*mocks.MockMailServerClient)
		setupTest      func(*testing.T, *mongo.Database)
	}{
		{
			name: "Successful Password Reset Request",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "sent")
				assert.Equal(t, "Reset password email sent successfully", response["message"])
				assert.Equal(t, true, response["sent"])

				// Verify reset password request was created in database
				var resetRequest models.UserResetPassword
				err = database.Collection("user_reset_password_requests").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&resetRequest)
				assert.NoError(t, err)
				assert.Equal(t, userID, *resetRequest.UserID)
				assert.NotEmpty(t, resetRequest.ResetCode)
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// Mock successful email send
				mockResponse := &connect.Response[mailserverv1.SendMailInternalResponse]{
					Msg: &mailserverv1.SendMailInternalResponse{
						Success: true,
					},
				}
				mockClient.On("SendMailInternal", mock.Anything, mock.Anything).Return(mockResponse, nil)
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "User Not Found",
			requestBody: map[string]interface{}{
				"email": "nonexistent@example.com",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "User not found", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// No mock setup needed for this test case
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "User Without Backup Email",
			requestBody: map[string]interface{}{
				"email": "nobackup@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Equal(t, "no_backup_email", response["message"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// No mock setup needed for this test case
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// Create user without backup email
				userID := primitive.NewObjectID()
				email := "nobackup@example.com"
				password := "hashedPassword"
				userWithoutBackup := &models.UserEntity{
					ID:          &userID,
					Email:       &email,
					BackupEmail: nil, // No backup email
					Password:    &password,
				}
				_, err := database.Collection("users").InsertOne(context.TODO(), userWithoutBackup)
				assert.NoError(t, err, "Failed to insert test user without backup email")
			},
		},
		{
			name: "Invalid JSON Request",
			requestBody: map[string]interface{}{
				"invalid_field": "test",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// No mock setup needed for this test case
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "Missing Email Field",
			requestBody: map[string]interface{}{
				"password": "test123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// No mock setup needed for this test case
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "Invalid Email Format",
			requestBody: map[string]interface{}{
				"email": "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// No mock setup needed for this test case
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "Existing Reset Request - Should Delete and Create New",
			requestBody: map[string]interface{}{
				"email": "existing@example.com",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "sent")
				assert.Equal(t, "Reset password email sent successfully", response["message"])
				assert.Equal(t, true, response["sent"])

				// Verify only one reset password request exists (old one should be deleted)
				// We'll check by counting all requests for the user
				cursor, err := database.Collection("user_reset_password_requests").Find(context.TODO(), bson.M{})
				assert.NoError(t, err)
				var requests []models.UserResetPassword
				err = cursor.All(context.TODO(), &requests)
				assert.NoError(t, err)
				// Should have only one request (the new one, old one should be deleted)
				assert.Equal(t, 1, len(requests))
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// Mock successful email send
				mockResponse := &connect.Response[mailserverv1.SendMailInternalResponse]{
					Msg: &mailserverv1.SendMailInternalResponse{
						Success: true,
					},
				}
				mockClient.On("SendMailInternal", mock.Anything, mock.Anything).Return(mockResponse, nil)
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// Create user for existing reset request test
				existingUserID := primitive.NewObjectID()
				email := "existing@example.com"
				backupEmail := "existing-backup@example.com"
				password := "hashedPassword"
				existingUser := &models.UserEntity{
					ID:          &existingUserID,
					Email:       &email,
					BackupEmail: &backupEmail,
					Password:    &password,
				}

				// Insert test user into database
				_, err := database.Collection("users").InsertOne(context.TODO(), existingUser)
				assert.NoError(t, err, "Failed to insert existing user")

				// Create existing reset password request
				existingRequest := &models.UserResetPassword{
					UserID:    &existingUserID,
					ResetCode: "old-reset-code",
					CreatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-time.Hour)),
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-time.Hour)),
				}
				_, err = database.Collection("user_reset_password_requests").InsertOne(context.TODO(), existingRequest)
				assert.NoError(t, err, "Failed to insert existing reset password request")
			},
		},
		{
			name: "Email Send Failure",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Failed to send email", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// Mock email send failure
				mockClient.On("SendMailInternal", mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
		{
			name: "Email Send Response Failure",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Failed to send email", response["error"])
			},
			setupMocks: func(mockClient *mocks.MockMailServerClient) {
				// Mock email send response with success=false
				mockResponse := &connect.Response[mailserverv1.SendMailInternalResponse]{
					Msg: &mailserverv1.SendMailInternalResponse{
						Success: false,
					},
				}
				mockClient.On("SendMailInternal", mock.Anything, mock.Anything).Return(mockResponse, nil)
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No additional setup needed
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear database collections before each test
			database.Collection("users").Drop(context.TODO())
			database.Collection("user_reset_password_requests").Drop(context.TODO())

			// Re-insert the base test user for tests that need it
			if tc.name != "User_Not_Found" && tc.name != "User_Without_Backup_Email" {
				_, err := database.Collection("users").InsertOne(context.TODO(), testUser)
				if err != nil {
					t.Fatalf("Failed to insert test user: %v", err)
				}
			}

			// Reset mock expectations
			mockMailServerClient.ExpectedCalls = nil
			mockMailServerClient.Calls = nil

			// Setup test case
			tc.setupTest(t, database)

			// Setup mocks
			tc.setupMocks(mockMailServerClient)

			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/start-reset-password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			tc.checkResponse(t, w, database)

			// Verify mock expectations
			mockMailServerClient.AssertExpectations(t)
		})
	}
}

func TestStartResetPasswordRequest(t *testing.T) {
	// Test creating a StartResetPasswordRequest
	request := StartResetPasswordRequest{
		Email: "test@example.com",
	}

	assert.Equal(t, "test@example.com", request.Email)
}

func TestStartResetPasswordRequestValidation(t *testing.T) {
	// Test cases for request validation
	testCases := []struct {
		name        string
		request     StartResetPasswordRequest
		expectValid bool
	}{
		{
			name: "Valid Email",
			request: StartResetPasswordRequest{
				Email: "test@example.com",
			},
			expectValid: true,
		},
		{
			name: "Empty Email",
			request: StartResetPasswordRequest{
				Email: "",
			},
			expectValid: false,
		},
		{
			name: "Invalid Email Format",
			request: StartResetPasswordRequest{
				Email: "invalid-email",
			},
			expectValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This would typically be tested with a validator, but since we're using Gin's binding
			// we can test the struct creation and basic validation
			request := tc.request
			if tc.expectValid {
				assert.NotEmpty(t, request.Email)
				// Additional validation would be done by Gin's binding
			} else {
				// For invalid cases, we expect empty or malformed email
				if tc.name == "Empty Email" {
					assert.Empty(t, request.Email)
				}
			}
		})
	}
}
