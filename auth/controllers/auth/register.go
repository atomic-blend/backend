package auth

import (
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/utils/jwt"
	"github.com/atomic-blend/backend/shared/utils/password"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Register creates a new user and returns tokens
// @Summary Register a new user
// @Description Register a new user with email and password
// @Accept json
// @Produce json
// @Param   request body RegisterRequest true "User registration data"
// @Success 201 {object} AuthResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.KeySet.Type != nil && *req.KeySet.Type != "age_v1" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key set type"})
		return
	}

	// check that the email domain is in the list of authorized domains from env variable
	authorizedDomains := os.Getenv("ACCOUNT_DOMAINS")
	if authorizedDomains == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ACCOUNT_DOMAINS not set"})
		return
	}
	authorizedDomainsList := strings.Split(authorizedDomains, ",")

	// check that the email domain is not in the list of restricted domains file
	restrictedDomains := os.Getenv("RESTRICTED_EMAILS")
	restrictedUsernames := strings.Split(restrictedDomains, ",")

	// check the username without domain or tags
	emailParts := strings.Split(req.Email, "@")
	emailUsername := emailParts[0]
	// remove tags from the username (everything after +)
	emailUsernameCleaned := strings.Split(emailUsername, "+")[0]

	if slices.Contains(restrictedUsernames, emailUsernameCleaned) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "restricted_email"})
		return
	}

	// Extract domain from email (email format is already validated by Gin binding)
	emailDomain := emailParts[1]
	if !slices.Contains(authorizedDomainsList, emailDomain) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Email domain is not authorized"})
		return
	}

	// Check if user already exists
	existingUser, err := c.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	// Hash the password
	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	// Find default user role
	defaultRole, err := c.userRoleRepo.FindOrCreate(ctx, "user")
	if err != nil {
		log.Error().Err(err).Msg("Failed to find default role")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find default role"})
		return
	}
	if defaultRole == nil {
		log.Error().Msg("Default role not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Default role not found"})
		return
	}

	// Create new user with default role
	user := &models.UserEntity{
		Email:       &req.Email,
		BackupEmail: req.BackupEmail,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Password:    &hashedPassword,
		KeySet:      req.KeySet,
		RoleIds:     []*primitive.ObjectID{defaultRole.ID},
	}

	// Save user to database
	newUser, err := c.userRepo.Create(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Populate user roles
	err = c.userRoleRepo.PopulateRoles(ctx, newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to populate user roles"})
		return
	}

	// Generate role list
	roles := make([]string, len(newUser.Roles))
	for i, role := range newUser.Roles {
		roles[i] = role.Name
	}

	// Generate tokens
	accessToken, err := jwt.GenerateToken(ctx, *newUser.ID, roles, jwt.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := jwt.GenerateToken(ctx, *newUser.ID, roles, jwt.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// For security reasons, remove the password from the response
	// Create a copy of the user without the password
	responseSafeUser := &models.UserEntity{
		ID:          newUser.ID,
		Email:       newUser.Email,
		FirstName:   newUser.FirstName,
		LastName:    newUser.LastName,
		BackupEmail: newUser.BackupEmail,
		KeySet:      newUser.KeySet,
		Roles:       newUser.Roles,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
	}

	// Return user and tokens
	ctx.JSON(http.StatusCreated, Response{
		User:         responseSafeUser,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
		ExpiresAt:    accessToken.ExpiresAt.Unix(),
	})
}
