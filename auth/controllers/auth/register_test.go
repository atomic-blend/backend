package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomic-blend/backend/auth/repositories"
	mailserver "github.com/atomic-blend/backend/shared/grpc/mail-server"
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

func TestRegister(t *testing.T) {
	// Set environment variable needed for JWT
	os.Setenv("SSO_SECRET", "test-secret-key")
	defer os.Unsetenv("SSO_SECRET")

	// Set environment variable for authorized domains
	os.Setenv("ACCOUNT_DOMAINS", "example.com,test.com,authorized.org")
	defer os.Unsetenv("ACCOUNT_DOMAINS")

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

	// Create user repository
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(database)
	mailServerClient, _ := mailserver.NewMailServerClient()

	// Create controller
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo, mailServerClient)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/register", authController.Register)

	// Create default user role before running tests
	roleID := primitive.NewObjectID()
	defaultRole := models.UserRoleEntity{
		ID:   &roleID,
		Name: "user",
	}
	_, err = database.Collection("user_roles").InsertOne(context.TODO(), defaultRole)
	if err != nil {
		t.Fatalf("Failed to insert default user role: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder, *mongo.Database)
		setupTest      func(*testing.T, *mongo.Database)
	}{
		{
			name: "Successful Registration",
			requestBody: map[string]interface{}{
				"email":    "newuser@example.com",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify tokens and basic user info
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotZero(t, response.ExpiresAt)
				assert.NotNil(t, response.User)
				assert.Equal(t, "newuser@example.com", *response.User.Email)

				// Verify role in response
				assert.NotNil(t, response.User.Roles)
				assert.Equal(t, 1, len(response.User.Roles))
				assert.Equal(t, roleID, *response.User.Roles[0].ID)
				assert.Equal(t, "user", response.User.Roles[0].Name)

				// Verify user in database
				var savedUser models.UserEntity
				err = database.Collection("users").FindOne(context.TODO(), bson.M{"email": "newuser@example.com"}).Decode(&savedUser)
				assert.NoError(t, err)

				// Verify password is hashed
				assert.NotNil(t, savedUser.Password)
				assert.NotEqual(t, "securePassword123", *savedUser.Password)

				// Verify KeySet
				assert.NotNil(t, savedUser.KeySet)
				assert.Equal(t, "encryptedUserKey123", savedUser.KeySet.UserKey)
				assert.Equal(t, "encryptedBackupKey123", savedUser.KeySet.BackupKey)
				assert.Equal(t, "salt123", savedUser.KeySet.Salt)
				assert.Equal(t, "mnemonicSalt123", savedUser.KeySet.MnemonicSalt)

				// Verify role assignment in database
				assert.NotNil(t, savedUser.RoleIds)
				assert.Equal(t, 1, len(savedUser.RoleIds))
				assert.Equal(t, roleID, *savedUser.RoleIds[0])
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed
			},
		},
		{
			name: "Email Already Exists",
			requestBody: map[string]interface{}{
				"email":    "existing@example.com",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Email is already registered", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// Create an existing user
				email := "existing@example.com"
				password := "hashedPassword" // In a real scenario this would be hashed
				_, err := database.Collection("users").InsertOne(context.TODO(), models.UserEntity{
					Email:    &email,
					Password: &password,
				})
				assert.NoError(t, err, "Failed to insert test user")
			},
		},
		{
			name: "Password Too Short",
			requestBody: map[string]interface{}{
				"email":    "valid@example.com",
				"password": "short",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Missing Email",
			requestBody: map[string]interface{}{
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Missing Password",
			requestBody: map[string]interface{}{
				"email": "valid@example.com",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Missing KeySet",
			requestBody: map[string]interface{}{
				"email":    "valid@example.com",
				"password": "securePassword123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Email Domain Not in Authorized Domains",
			requestBody: map[string]interface{}{
				"email":    "unauthorized@notallowed.com",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusForbidden,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Email domain is not authorized", response["error"])
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Email in Authorized Domains - Success",
			requestBody: map[string]interface{}{
				"email":    "authorized@test.com",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify tokens and basic user info
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotZero(t, response.ExpiresAt)
				assert.NotNil(t, response.User)
				assert.Equal(t, "authorized@test.com", *response.User.Email)

				// Verify user was created in database
				var savedUser models.UserEntity
				err = database.Collection("users").FindOne(context.TODO(), bson.M{"email": "authorized@test.com"}).Decode(&savedUser)
				assert.NoError(t, err)
				assert.Equal(t, "authorized@test.com", *savedUser.Email)
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Invalid Email Format - No @ Symbol",
			requestBody: map[string]interface{}{
				"email":    "invalidemail",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Contains(t, response["error"], "email")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Invalid Email Format - Multiple @ Symbols",
			requestBody: map[string]interface{}{
				"email":    "invalid@email@format.com",
				"password": "securePassword123",
				"keySet": map[string]interface{}{
					"userKey":      "encryptedUserKey123",
					"backupKey":    "encryptedBackupKey123",
					"salt":         "salt123",
					"mnemonicSalt": "mnemonicSalt123",
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, database *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Contains(t, response["error"], "email")
			},
			setupTest: func(t *testing.T, database *mongo.Database) {
				// No setup needed for this test case
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test case
			tc.setupTest(t, database)

			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
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

func TestRegister_AccountDomainsNotSet(t *testing.T) {
	// Set environment variable needed for JWT
	os.Setenv("SSO_SECRET", "test-secret-key")
	defer os.Unsetenv("SSO_SECRET")

	// Explicitly unset ACCOUNT_DOMAINS to test the error case
	os.Unsetenv("ACCOUNT_DOMAINS")

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

	// Create user repository
	userRepo := userrepo.NewUserRepository(database)
	userRoleRepo := userrolerepo.NewUserRoleRepository(database)
	resetPasswordRepo := repositories.NewUserResetPasswordRequestRepository(database)
	mailServerClient, _ := mailserver.NewMailServerClient()
	// Create controller
	authController := NewController(userRepo, userRoleRepo, resetPasswordRepo, mailServerClient	)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/register", authController.Register)

	// Create default user role before running tests
	roleID := primitive.NewObjectID()
	defaultRole := models.UserRoleEntity{
		ID:   &roleID,
		Name: "user",
	}
	_, err = database.Collection("user_roles").InsertOne(context.TODO(), defaultRole)
	if err != nil {
		t.Fatalf("Failed to insert default user role: %v", err)
	}

	// Test case for when ACCOUNT_DOMAINS is not set
	requestBody := map[string]interface{}{
		"email":    "test@example.com",
		"password": "securePassword123",
		"keySet": map[string]interface{}{
			"userKey":      "encryptedUserKey123",
			"backupKey":    "encryptedBackupKey123",
			"salt":         "salt123",
			"mnemonicSalt": "mnemonicSalt123",
		},
	}

	// Create request body
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Check response
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "ACCOUNT_DOMAINS not set", response["error"])
}
