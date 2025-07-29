package jwt

import (
	"context"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/atomic-blend/backend/auth/tests/utils/inmemorymongo"
	"github.com/atomic-blend/backend/auth/utils/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestDB(t *testing.T) func() {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	require.NoError(t, err)

	// Set up global database reference
	db.Database = client.Database("test_db")

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
		db.Database = nil
	}

	return cleanup
}

func TestGenerateToken(t *testing.T) {
	// Setup database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret")
	defer os.Setenv("SSO_SECRET", originalSecret)

	userID := primitive.NewObjectID()

	// Create a test gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	t.Run("should generate access token successfully", func(t *testing.T) {
		td, err := GenerateToken(ctx, userID, []string{
			"admin",
		}, AccessToken)

		assert.NoError(t, err)
		assert.NotEmpty(t, td.Token)
		assert.Equal(t, AccessToken, td.TokenType)
		assert.Equal(t, userID.Hex(), td.UserID)
		assert.Equal(t, []string{"admin"}, td.Roles)
		assert.True(t, td.ExpiresAt.After(time.Now()))
	})
}

func TestValidateToken(t *testing.T) {
	// Setup database
	cleanup := setupTestDB(t)
	defer cleanup()

	// Setup
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret")
	defer os.Setenv("SSO_SECRET", originalSecret)

	userID := primitive.NewObjectID()

	// Create a test gin context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	t.Run("should validate valid token successfully", func(t *testing.T) {
		td, err := GenerateToken(ctx, userID, []string{"admin"}, AccessToken)
		assert.NoError(t, err)

		claims, err := ValidateToken(td.Token, AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID.Hex(), (*claims)["user_id"])
		assert.Equal(t, string(AccessToken), (*claims)["type"])

		// Convert roles from interface{} to []string
		rolesInterface := (*claims)["roles"]
		roles, ok := rolesInterface.([]interface{})
		assert.True(t, ok)
		expectedRoles := make([]string, len(roles))
		for i, role := range roles {
			expectedRoles[i] = role.(string)
		}
		assert.Equal(t, []string{"admin"}, expectedRoles)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		claims, err := ValidateToken("invalid.token.string", AccessToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("should fail with wrong token type", func(t *testing.T) {
		td, err := GenerateToken(ctx, userID, []string{"admin"}, AccessToken)
		assert.NoError(t, err)

		claims, err := ValidateToken(td.Token, RefreshToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, "invalid token type", err.Error())
	})

	t.Run("should fail when SSO_SECRET not set", func(t *testing.T) {
		os.Setenv("SSO_SECRET", "")
		td, err := GenerateToken(ctx, userID, []string{"admin"}, AccessToken)
		assert.NoError(t, err) // Should still generate with default secret

		claims, err := ValidateToken(td.Token, AccessToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, "SSO_SECRET not set", err.Error())
	})
}

func TestGenerateJWKS(t *testing.T) {
	t.Run("should generate JWKS successfully", func(t *testing.T) {
		originalSecret := os.Getenv("SSO_SECRET")
		os.Setenv("SSO_SECRET", "test_secret")
		defer os.Setenv("SSO_SECRET", originalSecret)

		key, err := GenerateJWKS("HS256")
		assert.NoError(t, err)
		assert.NotNil(t, key)
	})

	t.Run("should fail when SSO_SECRET not set", func(t *testing.T) {
		os.Setenv("SSO_SECRET", "")
		key, err := GenerateJWKS("HS256")
		assert.Error(t, err)
		assert.Nil(t, key)
		assert.Equal(t, "SSO_SECRET not set", err.Error())
	})
}
