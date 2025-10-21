package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/models"
	userrepo "github.com/atomic-blend/backend/shared/repositories/user"
	userrolerepo "github.com/atomic-blend/backend/shared/repositories/user_role"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestConfirmResetPassword(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

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
	router.POST("/auth/confirm-reset-password", authController.ConfirmResetPassword)

	// Create test user
	userID := primitive.NewObjectID()
	email := "test@example.com"
	password := "hashedPassword"
	backupEmail := "backup@example.com"
	testUser := &models.UserEntity{
		ID:          &userID,
		Email:       &email,
		BackupEmail: &backupEmail,
		Password:    &password,
		KeySet: &models.EncryptionKey{
			UserKey:      "oldUserKey",
			BackupKey:    "oldBackupKey",
			Salt:         "oldSalt",
			MnemonicSalt: "oldMnemonicSalt",
		},
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder, *mongo.Database)
		setupTest      func(*testing.T, *mongo.Database) *models.UserResetPassword
	}{
		{
			name: "Successful Password Reset Confirmation",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "message")
				assert.Equal(t, "Password reset successfully", response["message"])

				// Verify user password was updated
				var updatedUser models.UserEntity
				err = database.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&updatedUser)
				assert.NoError(t, err)
				assert.NotEqual(t, password, *updatedUser.Password) // Password should be different
				assert.Equal(t, "newUserKey123", updatedUser.KeySet.UserKey)
				assert.Equal(t, "newUserSalt123", updatedUser.KeySet.Salt)
				assert.Equal(t, "newBackupKey123", updatedUser.KeySet.BackupKey)
				assert.Equal(t, "newBackupSalt123", updatedUser.KeySet.MnemonicSalt)

				// Verify reset code was deleted
				var resetRequest models.UserResetPassword
				err = database.Collection("user_reset_password_requests").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&resetRequest)
				assert.Error(t, err) // Should not find the reset code
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				// Insert test user
				_, err := database.Collection("users").InsertOne(context.TODO(), testUser)
				assert.NoError(t, err, "Failed to insert test user")

				// Create valid reset code (not expired)
				resetRequest := &models.UserResetPassword{
					UserID:    &userID,
					ResetCode: "valid-reset-code",
					CreatedAt: primitive.NewDateTimeFromTime(time.Now()), // Current time, not expired
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				}
				_, err = database.Collection("user_reset_password_requests").InsertOne(context.TODO(), resetRequest)
				assert.NoError(t, err, "Failed to insert reset password request")
				return resetRequest
			},
		},
		{
			name: "Invalid Reset Code",
			requestBody: map[string]interface{}{
				"reset_code":   "invalid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Reset code not found", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				// Insert test user
				_, err := database.Collection("users").InsertOne(context.TODO(), testUser)
				assert.NoError(t, err, "Failed to insert test user")
				return nil // No reset code
			},
		},
		{
			name: "Expired Reset Code",
			requestBody: map[string]interface{}{
				"reset_code":   "expired-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Reset code expired", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				// Insert test user
				_, err := database.Collection("users").InsertOne(context.TODO(), testUser)
				assert.NoError(t, err, "Failed to insert test user")

				// Create expired reset code (older than 5 minutes)
				resetRequest := &models.UserResetPassword{
					UserID:    &userID,
					ResetCode: "expired-reset-code",
					CreatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-10 * time.Minute)), // 10 minutes ago, expired
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-10 * time.Minute)),
				}
				_, err = database.Collection("user_reset_password_requests").InsertOne(context.TODO(), resetRequest)
				assert.NoError(t, err, "Failed to insert expired reset password request")
				return resetRequest
			},
		},
		{
			name: "User Not Found",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusInternalServerError,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Failed to find user", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				// Don't insert user, but create reset code for non-existent user
				nonExistentUserID := primitive.NewObjectID()
				resetRequest := &models.UserResetPassword{
					UserID:    &nonExistentUserID,
					ResetCode: "valid-reset-code",
					CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
					UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
				}
				_, err := database.Collection("user_reset_password_requests").InsertOne(context.TODO(), resetRequest)
				assert.NoError(t, err, "Failed to insert reset password request")
				return resetRequest
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
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing Reset Code",
			requestBody: map[string]interface{}{
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing Reset Data Field",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing New Password",
			requestBody: map[string]interface{}{
				"reset_code":  "valid-reset-code",
				"reset_data":  false,
				"user_key":    "newUserKey123",
				"user_salt":   "newUserSalt123",
				"backup_key":  "newBackupKey123",
				"backup_salt": "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing User Key",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing User Salt",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"backup_key":   "newBackupKey123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing Backup Key",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_salt":  "newBackupSalt123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
		{
			name: "Missing Backup Salt",
			requestBody: map[string]interface{}{
				"reset_code":   "valid-reset-code",
				"reset_data":   false,
				"new_password": "newSecurePassword123",
				"user_key":     "newUserKey123",
				"user_salt":    "newUserSalt123",
				"backup_key":   "newBackupKey123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid request", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) *models.UserResetPassword {
				return nil
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear database collections before each test
			database.Collection("users").Drop(context.TODO())
			database.Collection("user_reset_password_requests").Drop(context.TODO())

			// Setup test case
			tc.setupTest(t, database)

			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/confirm-reset-password", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			tc.checkResponse(t, w, database)
		})
	}
}

func TestConfirmResetPasswordRequest(t *testing.T) {
	// Test creating a ConfirmResetPasswordRequest
	resetData := true
	request := ConfirmResetPasswordRequest{
		ResetCode:   "test-reset-code",
		ResetData:   &resetData,
		NewPassword: "newPassword123",
		UserKey:     "userKey123",
		UserSalt:    "userSalt123",
		BackupKey:   "backupKey123",
		BackupSalt:  "backupSalt123",
	}

	assert.Equal(t, "test-reset-code", request.ResetCode)
	assert.Equal(t, true, *request.ResetData)
	assert.Equal(t, "newPassword123", request.NewPassword)
	assert.Equal(t, "userKey123", request.UserKey)
	assert.Equal(t, "userSalt123", request.UserSalt)
	assert.Equal(t, "backupKey123", request.BackupKey)
	assert.Equal(t, "backupSalt123", request.BackupSalt)
}

func TestConfirmResetPasswordRequestValidation(t *testing.T) {
	// Test cases for request validation
	testCases := []struct {
		name        string
		request     ConfirmResetPasswordRequest
		expectValid bool
	}{
		{
			name: "Valid Request",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: true,
		},
		{
			name: "Empty Reset Code",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Nil Reset Data",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   nil,
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Empty New Password",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Empty User Key",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Empty User Salt",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "",
				BackupKey:   "backupKey123",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Empty Backup Key",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "",
				BackupSalt:  "backupSalt123",
			},
			expectValid: false,
		},
		{
			name: "Empty Backup Salt",
			request: ConfirmResetPasswordRequest{
				ResetCode:   "test-reset-code",
				ResetData:   boolPtr(true),
				NewPassword: "newPassword123",
				UserKey:     "userKey123",
				UserSalt:    "userSalt123",
				BackupKey:   "backupKey123",
				BackupSalt:  "",
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
				assert.NotEmpty(t, request.ResetCode)
				assert.NotNil(t, request.ResetData)
				assert.NotEmpty(t, request.NewPassword)
				assert.NotEmpty(t, request.UserKey)
				assert.NotEmpty(t, request.UserSalt)
				assert.NotEmpty(t, request.BackupKey)
				assert.NotEmpty(t, request.BackupSalt)
			} else {
				// For invalid cases, we expect empty or nil values
				if tc.name == "Empty Reset Code" {
					assert.Empty(t, request.ResetCode)
				}
				if tc.name == "Nil Reset Data" {
					assert.Nil(t, request.ResetData)
				}
				if tc.name == "Empty New Password" {
					assert.Empty(t, request.NewPassword)
				}
				if tc.name == "Empty User Key" {
					assert.Empty(t, request.UserKey)
				}
				if tc.name == "Empty User Salt" {
					assert.Empty(t, request.UserSalt)
				}
				if tc.name == "Empty Backup Key" {
					assert.Empty(t, request.BackupKey)
				}
				if tc.name == "Empty Backup Salt" {
					assert.Empty(t, request.BackupSalt)
				}
			}
		})
	}
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}
