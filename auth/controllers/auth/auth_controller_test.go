package auth

import (
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/repositories"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewController(t *testing.T) {
	// Create mock repositories
	mockUserRepo := &repositories.UserRepository{}
	mockUserRoleRepo := &repositories.UserRoleRepository{}
	mockResetPasswordRepo := &repositories.UserResetPasswordRequestRepository{}

	// Create a new controller
	controller := NewController(mockUserRepo, mockUserRoleRepo, mockResetPasswordRepo)

	// Test that the controller was created successfully
	assert.NotNil(t, controller, "Controller should not be nil")

	// Test that the repositories were correctly assigned
	assert.Equal(t, mockUserRepo, controller.userRepo, "User repository should be correctly assigned")
	assert.Equal(t, mockUserRoleRepo, controller.userRoleRepo, "UserRole repository should be correctly assigned")

	// Test the controller type
	controllerType := reflect.TypeOf(controller)
	assert.Equal(t, "*auth.Controller", controllerType.String(), "Controller should be of type *auth.Controller")
}

func TestControllerStructure(t *testing.T) {
	// Create test cases for the structures
	t.Run("RegisterRequest structure", func(t *testing.T) {
		// Verify the RegisterRequest structure has the expected fields with expected tags
		registerRequestType := reflect.TypeOf(RegisterRequest{})

		// Check Email field
		emailField, found := registerRequestType.FieldByName("Email")
		assert.True(t, found, "Email field should exist")
		assert.Equal(t, "email", emailField.Tag.Get("json"))
		assert.Equal(t, "required,email", emailField.Tag.Get("binding"))

		// Check Password field
		passwordField, found := registerRequestType.FieldByName("Password")
		assert.True(t, found, "Password field should exist")
		assert.Equal(t, "password", passwordField.Tag.Get("json"))
		assert.Equal(t, "required,min=8", passwordField.Tag.Get("binding"))
	})

	t.Run("AuthResponse structure", func(t *testing.T) {
		// Verify the AuthResponse structure has the expected fields with expected tags
		authResponseType := reflect.TypeOf(Response{})

		// Check User field
		userField, found := authResponseType.FieldByName("User")
		assert.True(t, found, "User field should exist")
		assert.Equal(t, "*models.UserEntity", userField.Type.String())
		assert.Equal(t, "user", userField.Tag.Get("json"))

		// Check AccessToken field
		accessTokenField, found := authResponseType.FieldByName("AccessToken")
		assert.True(t, found, "AccessToken field should exist")
		assert.Equal(t, "string", accessTokenField.Type.String())
		assert.Equal(t, "accessToken", accessTokenField.Tag.Get("json"))

		// Check RefreshToken field
		refreshTokenField, found := authResponseType.FieldByName("RefreshToken")
		assert.True(t, found, "RefreshToken field should exist")
		assert.Equal(t, "string", refreshTokenField.Type.String())
		assert.Equal(t, "refreshToken", refreshTokenField.Tag.Get("json"))

		// Check ExpiresAt field
		expiresAtField, found := authResponseType.FieldByName("ExpiresAt")
		assert.True(t, found, "ExpiresAt field should exist")
		assert.Equal(t, "int64", expiresAtField.Type.String())
		assert.Equal(t, "expiresAt", expiresAtField.Tag.Get("json"))
	})

	t.Run("Controller structure", func(t *testing.T) {
		// Verify the Controller structure has the expected fields
		controllerType := reflect.TypeOf(Controller{})

		// Check userRepo field
		userRepoField, found := controllerType.FieldByName("userRepo")
		assert.True(t, found, "userRepo field should exist")
		assert.Equal(t, "repositories.UserRepositoryInterface", userRepoField.Type.String())

		// Check userRoleRepo field
		userRoleRepoField, found := controllerType.FieldByName("userRoleRepo")
		assert.True(t, found, "userRoleRepo field should exist")
		assert.Equal(t, "repositories.UserRoleRepositoryInterface", userRoleRepoField.Type.String())
	})
}

func TestRegisterRequest(t *testing.T) {
	// Test creating a RegisterRequest
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestAuthResponse(t *testing.T) {
	// Create a test user
	userID := primitive.NewObjectID()
	email := "test@example.com"
	user := &models.UserEntity{
		ID:    &userID,
		Email: &email,
	}

	// Test creating an AuthResponse
	resp := Response{
		User:         user,
		AccessToken:  "access-token-value",
		RefreshToken: "refresh-token-value",
		ExpiresAt:    1620000000,
	}

	assert.Equal(t, user, resp.User)
	assert.Equal(t, "access-token-value", resp.AccessToken)
	assert.Equal(t, "refresh-token-value", resp.RefreshToken)
	assert.Equal(t, int64(1620000000), resp.ExpiresAt)
}
