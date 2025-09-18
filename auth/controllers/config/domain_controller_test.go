package config

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewConfigController(t *testing.T) {
	controller := NewConfigController()
	assert.NotNil(t, controller)
	assert.IsType(t, &Controller{}, controller)
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

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

	// Create router
	router := gin.New()

	// Setup routes
	SetupRoutes(router, db)

	// Test that the route is properly set up by making a request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check that route was properly set up and returns 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response to verify it's the expected config response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "domains")
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create router
	router := gin.New()

	// Setup routes with mock
	SetupRoutesWithMock(router)

	// Test that the route is properly set up by making a request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check that route was properly set up and returns 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response to verify it's the expected config response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "domains")
}

func TestSetupConfigRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create controller
	controller := NewConfigController()

	// Create router
	router := gin.New()

	// Setup config routes directly
	setupConfigRoutes(router, controller)

	// Test that the route is properly set up by making a request
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(w, req)

	// Check that route was properly set up and returns 200
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response to verify it's the expected config response
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "domains")
}

func TestConfigRoutesGroup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create controller
	controller := NewConfigController()

	// Create router
	router := gin.New()

	// Setup config routes
	setupConfigRoutes(router, controller)

	// Test that the route group is properly configured
	// The route should be accessible at /config
	req, _ := http.NewRequest("GET", "/config", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 200 OK
	assert.Equal(t, http.StatusOK, w.Code)

	// Test that other routes in the group would be accessible
	// (though we only have one route currently)
	req2, _ := http.NewRequest("GET", "/config/nonexistent", nil)
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	// Should return 404 for non-existent routes
	assert.Equal(t, http.StatusNotFound, w2.Code)
}
