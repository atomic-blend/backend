package auth

import (
	"atomic_blend_api/repositories"
	"atomic_blend_api/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserAuthInfo contains the authenticated user information
type UserAuthInfo struct {
	UserID primitive.ObjectID
}

// Middleware verifies JWT tokens and adds user info to the context
// Can be applied to specific routes that require authentication
func Middleware() gin.HandlerFunc {
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

// requireRoleHandler checks if the authenticated user has the specified role
// It must be used after RequireAuth or AuthMiddleware
func requireRoleHandler(roleName string, userRepo *repositories.UserRepository, userRoleRepo *repositories.UserRoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authenticated user info
		authUser := GetAuthUser(c)
		if authUser == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get user details from database to check roles
		user, err := userRepo.FindByID(c, authUser.UserID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify user roles")
			if err.Error() == "user not found" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user roles"})
			}
			c.Abort()
			return
		}
		err = userRoleRepo.PopulateRoles(c, user)
		if err != nil {
			log.Error().Err(err).Msg("Failed to verify user roles")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify user roles"})
			c.Abort()
			return
		}
		log.Info().Msgf("User roles: %v", user.RoleIds)

		// Check if user has the required role
		hasRole := false
		for _, role := range user.Roles {
			if role != nil && role.Name == roleName {
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
