package mocks

import (
	"atomic_blend_api/models"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository provides a mock implementation of UserRepositoryInterface
type MockUserRepository struct {
	mock.Mock
}

// Create creates a new user
func (m *MockUserRepository) Create(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

// GetByID gets a user by ID
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.UserEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

// Update updates a user with the given ID
func (m *MockUserRepository) Update(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

// Delete deletes a user by ID
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// FindByEmail finds a user by email
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.UserEntity, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}

// FindByID finds a user by ID
func (m *MockUserRepository) FindByID(ctx *gin.Context, id primitive.ObjectID) (*models.UserEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserEntity), args.Error(1)
}
