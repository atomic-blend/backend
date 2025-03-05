package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/tests/utils/in_memory_mongo"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestRegister(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Start in-memory MongoDB server
	mongoServer, err := in_memory_mongo.CreateInMemoryMongoDB()
	if err != nil {
		t.Fatalf("Failed to create in-memory MongoDB: %v", err)
	}
	defer mongoServer.Stop()

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := in_memory_mongo.ConnectToInMemoryDB(mongoURI)
	if err != nil {
		t.Fatalf("Failed to connect to in-memory MongoDB: %v", err)
	}
	defer client.Disconnect(nil)

	// Get database reference
	db := client.Database("test_db")

	// Create user repository
	userRepo := repositories.NewUserRepository(db)

	// Create controller
	authController := NewController(userRepo)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/register", authController.Register)

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
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				// Check JSON response
				var response AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				// Verify tokens are present
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotZero(t, response.ExpiresAt)

				// Verify user details
				assert.NotNil(t, response.User)
				assert.Equal(t, "newuser@example.com", *response.User.Email)
				assert.Nil(t, response.User.Password, "Password should not be included in response")

				// Verify user was saved to database
				var savedUser models.UserEntity
				err = db.Collection("users").FindOne(nil, bson.M{"email": "newuser@example.com"}).Decode(&savedUser)
				assert.NoError(t, err)
				assert.NotNil(t, savedUser.Password, "Password should be stored (hashed)")
				assert.NotEqual(t, "securePassword123", *savedUser.Password, "Password should be hashed")
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Email Already Exists",
			requestBody: map[string]interface{}{
				"email":    "existing@example.com",
				"password": "securePassword123",
			},
			expectedStatus: http.StatusConflict,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Email is already registered", response["error"])
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// Create an existing user
				email := "existing@example.com"
				password := "hashedPassword" // In a real scenario this would be hashed
				_, err := db.Collection("users").InsertOne(nil, models.UserEntity{
					Email:    &email,
					Password: &password,
				})
				assert.NoError(t, err, "Failed to insert test user")
			},
		},
		{
			name: "Invalid Email Format",
			requestBody: map[string]interface{}{
				"email":    "invalidemail",
				"password": "securePassword123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Password Too Short",
			requestBody: map[string]interface{}{
				"email":    "valid@example.com",
				"password": "short",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Missing Email",
			requestBody: map[string]interface{}{
				"password": "securePassword123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// No setup needed for this test case
			},
		},
		{
			name: "Missing Password",
			requestBody: map[string]interface{}{
				"email": "valid@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder, db *mongo.Database) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
			setupTest: func(t *testing.T, db *mongo.Database) {
				// No setup needed for this test case
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup test case
			tc.setupTest(t, db)

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
			tc.checkResponse(t, w, db)
		})
	}
}
