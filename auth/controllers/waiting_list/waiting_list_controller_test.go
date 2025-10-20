package waitinglist

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	t.Run("should create new waiting list controller", func(t *testing.T) {
		// Setup
		mockRepo := new(mocks.MockWaitingListRepository)

		// Act
		controller := NewController(mockRepo)

		// Assert
		assert.NotNil(t, controller, "Controller should not be nil")
		assert.IsType(t, &Controller{}, controller)
		assert.Equal(t, mockRepo, controller.waitingListRepo, "Repository should be set correctly")
	})
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should setup waiting list routes with real database", func(t *testing.T) {
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
		req, _ := http.NewRequest("POST", "/waiting-list", nil)
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check that route was properly set up and returns 400 (bad request due to missing body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Parse response to verify it's the expected error response
		var response map[string]string
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})

	t.Run("should setup waiting list routes with mock repository", func(t *testing.T) {
		// Create mock repository
		mockRepo := new(mocks.MockWaitingListRepository)

		// Create controller
		controller := NewController(mockRepo)

		// Create router
		router := gin.New()
		waitingListGroup := router.Group("/waiting-list")
		{
			waitingListGroup.POST("", controller.JoinWaitingList)
		}

		// Test that the route is properly set up by making a request
		req, _ := http.NewRequest("POST", "/waiting-list", nil)
		w := httptest.NewRecorder()

		// Perform request
		router.ServeHTTP(w, req)

		// Check that route was properly set up and returns 400 (bad request due to missing body)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Parse response to verify it's the expected error response
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "error")
	})
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock repository
	mockRepo := new(mocks.MockWaitingListRepository)

	// Create controller
	controller := NewController(mockRepo)

	// Create router
	router := gin.New()
	waitingListGroup := router.Group("/waiting-list")
	{
		waitingListGroup.POST("", controller.JoinWaitingList)
	}

	// Test that the route group is properly configured
	// The route should be accessible at /waiting-list
	req, _ := http.NewRequest("POST", "/waiting-list", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should return 400 Bad Request for missing body
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test that other routes in the group would be accessible
	// (though we only have one route currently)
	req2, _ := http.NewRequest("GET", "/waiting-list/nonexistent", nil)
	w2 := httptest.NewRecorder()

	router.ServeHTTP(w2, req2)

	// Should return 404 for non-existent routes
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestControllerIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should handle route setup and basic functionality", func(t *testing.T) {
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

		// Test that the route group is properly configured
		routes := router.Routes()
		expectedRoutes := map[string]bool{
			"POST /waiting-list": false,
		}

		// Check that all expected routes are set up
		for _, route := range routes {
			routeKey := route.Method + " " + route.Path
			if _, exists := expectedRoutes[routeKey]; exists {
				expectedRoutes[routeKey] = true
			}
		}

		// Verify all expected routes were found
		for routeKey, found := range expectedRoutes {
			assert.True(t, found, "Expected route not found: %s", routeKey)
		}
	})
}
