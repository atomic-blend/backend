package auth

import (
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/utils/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RefreshToken handles token refresh requests
// @Summary Refresh access token
// @Description Generate new access and refresh tokens using a valid refresh token
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} AuthResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (c *Controller) RefreshToken(ctx *gin.Context) {
	// Get token from Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// For the test environment, handle the mock token differently
	// This is to make the test pass
	var userID primitive.ObjectID
	var err error

	// Validate refresh token normally
	claims, err := jwt.ValidateToken(tokenString, jwt.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract user ID from claims
	userIDStr, ok := (*claims)["user_id"].(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, err = primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}
	// Find user to ensure they still exist and get their current data
	user, err := c.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Populate user roles - we need to return immediately if there's an error
	if err = c.userRoleRepo.PopulateRoles(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate user roles"})
		return
	}

	// Generate new tokens
	accessToken, err := jwt.GenerateToken(ctx, userID, jwt.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(ctx, userID, jwt.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Return user and new tokens
	responseSafeUser := &models.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		KeySet:   user.KeySet,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, Response{
		User:         responseSafeUser,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt.Unix(),
	})
}
