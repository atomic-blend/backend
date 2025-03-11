package mocks

import (
	"atomic_blend_api/models"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)





// MockTaskRepository provides a mock implementation of TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

// GetAll gets all tasks
func (m *MockTaskRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TaskEntity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TaskEntity), args.Error(1)
}

// GetByID gets a task by ID
func (m *MockTaskRepository) GetByID(ctx context.Context, id string) (*models.TaskEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// Create creates a new task
func (m *MockTaskRepository) Create(ctx context.Context, task *models.TaskEntity) (*models.TaskEntity, error) {
	args := m.Called(ctx, task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// Update updates a task with the given ID
func (m *MockTaskRepository) Update(ctx context.Context, id string, task *models.TaskEntity) (*models.TaskEntity, error) {
	args := m.Called(ctx, id, task)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// Delete deletes a task by ID
func (m *MockTaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
