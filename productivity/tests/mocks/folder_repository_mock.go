package mocks

import (
	"atomic-blend/backend/productivity/models"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockFolderRepository provides a mock implementation of FolderRepositoryInterface
type MockFolderRepository struct {
	mock.Mock
}

// GetAll retrieves all folders for a user
func (m *MockFolderRepository) GetAll(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Folder), args.Error(1)
}

// Create creates a new folder
func (m *MockFolderRepository) Create(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	args := m.Called(ctx, folder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Folder), args.Error(1)
}

// Update modifies an existing folder
func (m *MockFolderRepository) Update(ctx context.Context, id primitive.ObjectID, folder *models.Folder) (*models.Folder, error) {
	args := m.Called(ctx, id, folder)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Folder), args.Error(1)
}

// Delete removes a folder
func (m *MockFolderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
