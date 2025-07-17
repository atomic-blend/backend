package userrole

import (
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewUserRoleController(t *testing.T) {
	// Setup
	mockRepo := new(mocks.MockUserRoleRepository)

	t.Run("should create new user role controller", func(t *testing.T) {
		// Act
		controller := NewUserRoleController(mockRepo)

		// Assert
		assert.NotNil(t, controller, "Controller should not be nil")
		assert.Equal(t, mockRepo, controller.userRoleRepo, "Repository should be set correctly")
	})
}

func TestSetupRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockUserRoleRepository)
	controller := NewUserRoleController(mockRepo)
	router := gin.New()
	group := router.Group("/admin")

	t.Run("should setup all user role routes", func(t *testing.T) {
		// Act
		controller.SetupRoutes(group)

		// Assert
		routes := router.Routes()
		expectedRoutes := map[string]bool{
			"GET /admin/user-roles":        false,
			"GET /admin/user-roles/:id":    false,
			"POST /admin/user-roles":       false,
			"PUT /admin/user-roles/:id":    false,
			"DELETE /admin/user-roles/:id": false,
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

