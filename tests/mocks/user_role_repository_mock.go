package mocks

import (
	"atomic_blend_api/models"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRoleRepository provides a mock implementation of UserRoleRepositoryInterface
type MockUserRoleRepository struct {
	mock.Mock
}

func (m *MockUserRoleRepository) Create(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UserRoleEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) GetAll(ctx context.Context) ([]*models.UserRoleEntity, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) GetByName(ctx context.Context, name string) (*models.UserRoleEntity, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) Update(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRoleRepository) FindOrCreate(ctx context.Context, roleName string) (*models.UserRoleEntity, error) {
	args := m.Called(ctx, roleName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserRoleEntity), args.Error(1)
}

func (m *MockUserRoleRepository) PopulateRoles(ctx context.Context, user *models.UserEntity) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}