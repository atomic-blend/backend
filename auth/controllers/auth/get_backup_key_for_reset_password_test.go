package auth

import (
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/repositories"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/password"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetBackupKeyForResetPassword(t *testing.T) {
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
	db := client.Database("test_db")

	// Create repositories
	userRepo := repositories.NewUserRepository(db)
	userRoleRepo := repositories.NewUserRoleRepository(db)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(db)

	// Create controller
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/get-backup-key", authController.GetBackupKeyForResetPassword)

	// Create a test role first
	roleID := primitive.NewObjectID()
	userRole := models.UserRoleEntity{
		ID:   &roleID,
		Name: "user",
	}
	_, err = db.Collection("user_roles").InsertOne(context.TODO(), userRole)
	if err != nil {
		t.Fatalf("Failed to insert test user role: %v", err)
	}

	// Create a test user with role reference
	hashedPassword, _ := password.HashPassword("testPassword123")
	testUserID := primitive.NewObjectID()
	keySet := models.EncryptionKey{
		UserKey:      "testUserKey123",
		BackupKey:    "testBackupKey123",
		Salt:         "testSalt123",
		MnemonicSalt: "testMnemonicSalt123",
	}
	testUser := models.UserEntity{
		ID:       &testUserID,
		Email:    stringPtr("test@example.com"),
		Password: &hashedPassword,
		KeySet:   &keySet,
		RoleIds:  []*primitive.ObjectID{&roleID},
	}

	// Insert test user into database
	_, err = db.Collection("users").InsertOne(context.TODO(), testUser)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Create a valid reset password request
	validResetCode := "valid-reset-code"
	resetPasswordRequest := models.UserResetPassword{
		UserID:    testUser.ID,
		ResetCode: validResetCode,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Insert reset password request into database
	_, err = db.Collection("user_reset_password_requests").InsertOne(context.TODO(), resetPasswordRequest)
	if err != nil {
		t.Fatalf("Failed to insert reset password request: %v", err)
	}

	// create another user for expired reset code test
	user2ID := primitive.NewObjectID()
	keySet2 := models.EncryptionKey{
		UserKey:      "testUserKey456",
		BackupKey:    "testBackupKey456",
		Salt:         "testSalt456",
		MnemonicSalt: "testMnemonicSalt456",
	}
	user2 := models.UserEntity{
		ID:       &user2ID,
		Email:    stringPtr("test2@example.com"),
		Password: &hashedPassword,
		KeySet:   &keySet2,
		RoleIds:  []*primitive.ObjectID{&roleID},
	}

	// Insert second test user into database
	_, err = db.Collection("users").InsertOne(context.TODO(), user2)
	if err != nil {
		t.Fatalf("Failed to insert second test user: %v", err)
	}

	// Create an expired reset password request
	expiredResetCode := "expired-reset-code"
	expiredResetPasswordRequest := models.UserResetPassword{
		UserID:    user2.ID,
		ResetCode: expiredResetCode,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-6 * time.Minute)), // 6 minutes ago (expired)
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now().Add(-6 * time.Minute)),
	}

	// Insert expired reset password request into database
	_, err = db.Collection("user_reset_password_requests").InsertOne(context.TODO(), expiredResetPasswordRequest)
	if err != nil {
		t.Fatalf("Failed to insert expired reset password request: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successful Retrieval of Backup Key",
			requestBody: map[string]interface{}{
				"reset_code": validResetCode,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Check if backup_key and backup_salt are present in the response
				assert.Equal(t, "testBackupKey123", response["backup_key"])
				assert.Equal(t, "testMnemonicSalt123", response["backup_salt"])
			},
		},
		{
			name:        "Invalid Request Format",
			requestBody: map[string]interface{}{
				// Missing reset_code
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "Reset Code Not Found",
			requestBody: map[string]interface{}{
				"reset_code": "non-existent-code",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Reset code not found", response["error"])
			},
		},
		{
			name: "Expired Reset Code",
			requestBody: map[string]interface{}{
				"reset_code": expiredResetCode,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Reset code expired", response["error"])
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert request body to JSON
			jsonBody, err := json.Marshal(tc.requestBody)
			assert.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest("POST", "/auth/get-backup-key", bytes.NewBuffer(jsonBody))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			if tc.checkResponse != nil {
				tc.checkResponse(t, w)
			}
		})
	}
}
