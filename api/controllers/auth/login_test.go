package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/repositories"
	"atomic_blend_api/tests/utils/in_memory_mongo"
	"atomic_blend_api/utils/password"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
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
	userRoleRepo := repositories.NewUserRoleRepository(db)

	// Create controller
	authController := NewController(userRepo, userRoleRepo)

	// Create a test router
	router := gin.Default()
	router.POST("/auth/login", authController.Login)

	// Create a test role first
	roleID := primitive.NewObjectID()
	userRole := models.UserRoleEntity{
		ID:   &roleID,
		Name: "user",
	}
	_, err = db.Collection("user_roles").InsertOne(nil, userRole)
	if err != nil {
		t.Fatalf("Failed to insert test user role: %v", err)
	}

	// Create a test user with role reference
	hashedPassword, _ := password.HashPassword("testPassword123")
	testUserID := primitive.NewObjectID()
	keySalt := "0123456789abcdef0123456789abcdef" // Test key salt
	testUser := models.UserEntity{
		ID:       &testUserID,
		Email:    stringPtr("test@example.com"),
		Password: &hashedPassword,
		KeySalt:  &keySalt,
		RoleIds:  []*primitive.ObjectID{&roleID}, // Add role reference
	}

	// Insert test user into database
	_, err = db.Collection("users").InsertOne(nil, testUser)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	// Test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "Successful Login",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "testPassword123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response AuthResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.NotNil(t, response.User)
				assert.Equal(t, "test@example.com", *response.User.Email)
				assert.Nil(t, response.User.Password) // Password should not be returned
				assert.NotNil(t, response.User.KeySalt)
				assert.Equal(t, keySalt, *response.User.KeySalt)

				// Verify roles are populated
				assert.NotNil(t, response.User.Roles)
				assert.Equal(t, 1, len(response.User.Roles))
				assert.Equal(t, roleID, *response.User.Roles[0].ID)
				assert.Equal(t, "user", response.User.Roles[0].Name)
			},
		},
		{
			name: "Invalid Password",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "wrongPassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid email or password", response["error"])
			},
		},
		{
			name: "User Not Found",
			requestBody: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "testPassword123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
				assert.Equal(t, "Invalid email or password", response["error"])
			},
		},
		{
			name: "Missing Email",
			requestBody: map[string]interface{}{
				"password": "testPassword123",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
		{
			name: "Missing Password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "error")
			},
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request body
			jsonBody, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Check response
			tc.checkResponse(t, w)
		})
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
