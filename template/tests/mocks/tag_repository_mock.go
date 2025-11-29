package mocks

import (
	"context"

	"github.com/atomic-blend/backend/mail/models"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTagRepository provides a mock implementation of TagRepository
type MockTagRepository struct {
	mock.Mock
}

// GetAll gets all tags
func (m *MockTagRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Tag, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Tag), args.Error(1)
}

// GetByID gets a tag by ID
func (m *MockTagRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

// Create creates a new tag
func (m *MockTagRepository) Create(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

// Update updates a tag with the given ID
func (m *MockTagRepository) Update(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	args := m.Called(ctx, tag)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Tag), args.Error(1)
}

// Delete deletes a tag by ID
func (m *MockTagRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DeleteByUserID deletes all tags for a specific user
func (m *MockTagRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
