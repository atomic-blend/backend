package health

import (
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetHealth(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		disconnectDB   bool
		expectedStatus int
	}{
		{
			name:           "Both API and DB are healthy",
			disconnectDB:   false,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "API is up but DB is down",
			disconnectDB:   true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			// Get database reference
			db := client.Database("test_db")

			// If we need to simulate DB down scenario, we'll disconnect the client
			if tc.disconnectDB {
				// Force disconnect to simulate DB down
				client.Disconnect(context.Background())
			} else {
				defer client.Disconnect(context.Background())
			}

			// Create controller
			healthController := NewHealthController(db)

			// Create test router
			router := gin.Default()
			router.GET("/health", healthController.GetHealth)

			// Create request
			req, _ := http.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Parse response body
			var response Response
			err = json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check the response values
			assert.Equal(t, true, response.Up, "API should always be up")

			if tc.disconnectDB {
				assert.Equal(t, false, response.DB, "DB should be down")
			} else {
				assert.Equal(t, true, response.DB, "DB should be up")
			}
		})
	}
}

// Additional test for SetupRoutes
func TestSetupRoutes(t *testing.T) {
	// Create a new router
	router := gin.Default()

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

	// Setup routes
	SetupRoutes(router, db)

	// Create request to verify route is setup
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check that route was properly set up and returns 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response to verify it's the expected health response
	var response Response
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, true, response.Up)
}
