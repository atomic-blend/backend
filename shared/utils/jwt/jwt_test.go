package jwt

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret")
	defer os.Setenv("SSO_SECRET", originalSecret)

	userID := primitive.NewObjectID()

	t.Run("should generate access token successfully", func(t *testing.T) {
		td, err := GenerateToken(userID, AccessToken)

		assert.NoError(t, err)
		assert.NotEmpty(t, td.Token)
		assert.Equal(t, AccessToken, td.TokenType)
		assert.Equal(t, userID.Hex(), td.UserID)
		assert.True(t, td.ExpiresAt.After(time.Now()))
	})
}

func TestValidateToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret")
	defer os.Setenv("SSO_SECRET", originalSecret)

	userID := primitive.NewObjectID()

	t.Run("should validate valid token successfully", func(t *testing.T) {
		td, err := GenerateToken(userID, AccessToken)
		assert.NoError(t, err)

		claims, err := ValidateToken(td.Token, AccessToken)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, userID.Hex(), *claims.UserID)
		assert.Equal(t, string(AccessToken), *claims.Type)
	})

	t.Run("should fail with invalid token", func(t *testing.T) {
		claims, err := ValidateToken("invalid.token.string", AccessToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("should fail with wrong token type", func(t *testing.T) {
		td, err := GenerateToken(userID, AccessToken)
		assert.NoError(t, err)

		claims, err := ValidateToken(td.Token, RefreshToken)
		assert.Error(t, err)
		assert.Nil(t, claims)
		assert.Equal(t, "invalid token type", err.Error())
	})

	t.Run("should fail when SSO_SECRET not set", func(t *testing.T) {
		os.Setenv("SSO_SECRET", "")
		td, err := GenerateToken(userID, AccessToken)
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
