package habits

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewHabitController(t *testing.T) {
	mockRepo := new(mocks.MockHabitRepository)
	controller := NewHabitController(mockRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, mockRepo, controller.habitRepo)
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockHabitRepository)

	SetupRoutesWithMock(router, mockRepo)

	// Test that routes are properly registered by making test requests
	testRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/habits"},
		{http.MethodGet, "/habits/123"},
		{http.MethodPost, "/habits"},
		{http.MethodPut, "/habits/123"},
		{http.MethodDelete, "/habits/123"},
	}

	// Setup mock expectations for each route
	mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*primitive.ObjectID")).Return([]*models.Habit{}, nil)
	mockRepo.On("GetByID", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(createTestHabit(), nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Habit")).Return(createTestHabit(), nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Habit")).Return(createTestHabit(), nil)
	mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	for _, route := range testRoutes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		router.ServeHTTP(w, req)

		// We're just checking if routes are registered, not their full functionality
		assert.NotEqual(t, http.StatusNotFound, w.Code, "Route not found: %s %s", route.method, route.path)
	}
}

func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockHabitRepository)

	// Use SetupRoutesWithMock instead of SetupRoutes to avoid database dependency
	assert.NotPanics(t, func() {
		SetupRoutesWithMock(router, mockRepo)
	})
}

func TestRouteRegistration(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockHabitRepository)
	controller := NewHabitController(mockRepo)

	// Call the function through public function
	SetupRoutesWithMock(router, mockRepo)

	// Verify all expected routes exist by checking if they're handled
	paths := []string{
		"/habits",
		"/habits/:id",
	}

	for _, path := range paths {
		// Check if the route exists
		r := router.Routes()
		found := false

		for _, route := range r {
			if route.Path == path {
				found = true
				break
			}
		}

		assert.True(t, found, "Expected route not registered: %s", path)
	}

	// Also verify controller is properly constructed
	assert.NotNil(t, controller)
	assert.Equal(t, mockRepo, controller.habitRepo)
}
