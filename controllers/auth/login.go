package auth

import (
	"atomic_blend_api/models"
	"atomic_blend_api/utils/jwt"
	"atomic_blend_api/utils/password"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest represents the structure for login request data
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns tokens
// @Summary Login user
// @Description Authenticate a user with email and password and return tokens
// @Accept json
// @Produce json
// @Param   request body LoginRequest true "User login data"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by email
	user, err := c.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	if user.Password == nil || !password.CheckPassword(req.Password, *user.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Populate user roles
	err = c.userRoleRepo.PopulateRoles(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate user roles"})
		return
	}

	// Generate tokens
	accessToken, err := jwt.GenerateToken(*user.ID, jwt.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(*user.ID, jwt.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// For security reasons, remove the password from the response
	responseSafeUser := &models.UserEntity{
		ID:        user.ID,
		Email:     user.Email,
		KeySet:   user.KeySet,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// Return user and tokens
	ctx.JSON(http.StatusOK, Response{
		User:         responseSafeUser,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt.Unix(),
	})
}
