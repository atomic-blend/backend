package users

import (
	"atomic_blend_api/repositories"
	"atomic_blend_api/tests/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserController(t *testing.T) {
	// Create mock repositories
	mockUserRepo := new(mocks.MockUserRepository)
	mockUserRoleRepo := new(mocks.MockUserRoleRepository)

	// Create controller
	controller := NewUserController(mockUserRepo, mockUserRoleRepo)

	// Assert controller properties
	assert.NotNil(t, controller)
	assert.Equal(t, mockUserRepo, controller.userRepo)
	assert.Equal(t, mockUserRoleRepo, controller.userRoleRepo)
}

func TestUserControllerImplementsInterfaces(t *testing.T) {
	// This is just a compile-time check to ensure the interfaces match
	var _ repositories.UserRepositoryInterface = &mocks.MockUserRepository{}
	var _ repositories.UserRoleRepositoryInterface = &mocks.MockUserRoleRepository{}
}
