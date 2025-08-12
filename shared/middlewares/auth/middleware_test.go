package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomic-blend/backend/shared/utils/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthMiddleware(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_auth_middleware")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header is required")
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "InvalidFormat")
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header must be in format: Bearer {token}")
	})

	t.Run("Invalid Token", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer invalid-token")
		c.Request = req

		// Execute middleware
		Middleware()(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Setup
		userID := primitive.NewObjectID()
		tokenDetails, err := jwt.GenerateToken(userID, jwt.AccessToken)
		assert.NoError(t, err, "Token generation should not fail")
		assert.NotEmpty(t, tokenDetails.Token, "Token should not be empty")

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+tokenDetails.Token)
		c.Request = req

		// Create a new engine to use with the middleware
		router := gin.New()
		router.Use(Middleware())
		router.GET("/", func(c *gin.Context) {
			// Assert user was set in context
			authUser := GetAuthUser(c)
			assert.NotNil(t, authUser)
			assert.Equal(t, userID, authUser.UserID)
			c.Status(http.StatusOK)
		})

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetAuthUser(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("Auth User Not Set", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.Nil(t, authUser)
	})

	t.Run("Auth User Set", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		userID := primitive.NewObjectID()
		c.Set("authUser", &UserAuthInfo{UserID: userID})

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.NotNil(t, authUser)
		assert.Equal(t, userID, authUser.UserID)
	})

	t.Run("Invalid Auth User Type", func(t *testing.T) {
		// Setup
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Set invalid type
		c.Set("authUser", "not-a-user-auth-info")

		// Execute
		authUser := GetAuthUser(c)

		// Assert
		assert.Nil(t, authUser)
	})
}

// Helper function to wrap our mock repositories for testing
func mockRoleHandler(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user info
		authUser := GetAuthUser(c)
		if authUser == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Check if Claims is nil
		if authUser.Claims == nil || authUser.Claims.Roles == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		roles := *authUser.Claims.Roles

		// Check if user has the required role
		hasRole := false
		for _, role := range roles {
			if role == roleName {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// User has the required role, proceed
		c.Next()
	}
}

func TestRequireRoleHandler(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("No Auth User", func(t *testing.T) {
		// Setup mocks
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Execute our test wrapper function
		mockRoleHandler("admin")(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "Authentication required")
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Setup mocks
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Set authUser with nil Claims to simulate user not found
		c.Set("authUser", &UserAuthInfo{UserID: userID, Claims: nil})

		// Execute our test wrapper function
		mockRoleHandler("admin")(c)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("User Has Required Role", func(t *testing.T) {
		// Setup mocks
		userID := primitive.NewObjectID()

		// Create a new router to test the middleware chain
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Set auth user in context with proper Claims
			roles := []string{"admin"}
			c.Set("authUser", &UserAuthInfo{
				UserID: userID,
				Claims: &jwt.CustomClaims{
					Roles: &roles,
				},
			})
		})
		router.Use(mockRoleHandler("admin"))
		router.GET("/", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create request
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("User Doesn't Have Required Role", func(t *testing.T) {
		// Setup mocks

		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		// Set authUser with user role (not admin)
		roles := []string{"user"}
		c.Set("authUser", &UserAuthInfo{
			UserID: userID,
			Claims: &jwt.CustomClaims{
				Roles: &roles,
			},
		})

		// Execute our test wrapper function
		mockRoleHandler("admin")(c)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})
}

func TestOptionalAuth(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_optional_auth")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("No Authorization Header", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Auth user should not be set
			authUser := GetAuthUser(c)
			assert.Nil(t, authUser)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		// Execute
		router.ServeHTTP(w, req)

		// Assert next was called and no errors
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Setup
		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Auth user should not be set
			authUser := GetAuthUser(c)
			assert.Nil(t, authUser)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "InvalidFormat")

		// Execute
		router.ServeHTTP(w, req)

		// Assert next was called and no errors
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Valid Token", func(t *testing.T) {
		// Setup
		userID := primitive.NewObjectID()
		tokenDetails, err := jwt.GenerateToken(userID, jwt.AccessToken)
		assert.NoError(t, err, "Token generation should not fail")

		router := gin.New()
		router.Use(OptionalAuth())
		router.GET("/", func(c *gin.Context) {
			// Assert user was set in context
			authUser := GetAuthUser(c)
			if authUser == nil {
				t.Fatal("Auth user should not be nil")
			}
			assert.Equal(t, userID, authUser.UserID)
			c.Status(http.StatusOK)
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+tokenDetails.Token)

		// Execute
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
