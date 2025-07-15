package admin

import (
	"context"
	"testing"

	"auth/tests/utils/inmemorymongo"
	"auth/utils/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup in-memory MongoDB
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)
	defer mongoServer.Stop()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)
	defer client.Disconnect(context.Background())

	// Set global database for auth middleware
	database := client.Database("test_db")
	db.Database = database

	t.Run("should setup admin routes with middleware", func(t *testing.T) {
		// Act
		SetupRoutes(router, database)

		// Assert
		routes := router.Routes()
		hasAdminRoute := false
		for _, route := range routes {
			if route.Path == "/admin/user-roles" {
				hasAdminRoute = true
				break
			}
		}
		assert.True(t, hasAdminRoute, "Admin routes should be set up")
	})
}

func TestNewAdminController(t *testing.T) {
	t.Run("should create new admin controller", func(t *testing.T) {
		// Act
		controller := NewAdminController()

		// Assert
		assert.NotNil(t, controller, "Controller should not be nil")
	})
}
