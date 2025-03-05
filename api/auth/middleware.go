package auth

import (
	"atomic_blend_api/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserAuthInfo contains the authenticated user information
type UserAuthInfo struct {
	UserID primitive.ObjectID
}

// AuthMiddleware verifies JWT tokens and adds user info to the context
// Can be applied to specific routes that require authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must be in format: Bearer {token}"})
			c.Abort()
			return
		}

		// Extract and validate the token
		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString, jwt.AccessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from token claims
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: missing user_id claim"})
			c.Abort()
			return
		}

		// Convert string user ID to ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		// Set user info in context for use in subsequent handlers
		c.Set("authUser", &UserAuthInfo{
			UserID: userID,
		})

		c.Next()
	}
}

// GetAuthUser retrieves the authenticated user info from the Gin context
// Use this in your handlers after applying the AuthMiddleware
func GetAuthUser(c *gin.Context) *UserAuthInfo {
	userValue, exists := c.Get("authUser")
	if !exists {
		return nil
	}

	user, ok := userValue.(*UserAuthInfo)
	if !ok {
		return nil
	}

	return user
}

// RequireAuth applies the auth middleware to a specific route group
// Example usage: RequireAuth(router.Group("/protected"))
func RequireAuth(group *gin.RouterGroup) *gin.RouterGroup {
	group.Use(AuthMiddleware())
	return group
}

// OptionalAuth middleware that doesn't abort if auth fails
// Useful for routes that work with different behavior for logged-in vs anonymous users
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Just continue without auth
			c.Next()
			return
		}

		// Check if the Authorization header has the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Just continue without auth
			c.Next()
			return
		}

		// Extract and validate the token
		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString, jwt.AccessToken)
		if err != nil {
			// Just continue without auth
			c.Next()
			return
		}

		// Extract user ID from token claims
		userIDStr, ok := (*claims)["user_id"].(string)
		if !ok {
			// Just continue without auth
			c.Next()
			return
		}

		// Convert string user ID to ObjectID
		userID, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			// Just continue without auth
			c.Next()
			return
		}

		// Set user info in context for use in subsequent handlers
		c.Set("authUser", &UserAuthInfo{
			UserID: userID,
		})

		c.Next()
	}
}
