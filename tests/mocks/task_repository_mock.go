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

// AddTimeEntry adds a time entry to a task
func (m *MockTaskRepository) AddTimeEntry(ctx context.Context, taskID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error) {
	args := m.Called(ctx, taskID, timeEntry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// RemoveTimeEntry removes a time entry from a task
func (m *MockTaskRepository) RemoveTimeEntry(ctx context.Context, taskID string, timeEntryID string) (*models.TaskEntity, error) {
	args := m.Called(ctx, taskID, timeEntryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// UpdateTimeEntry updates a time entry in a task
func (m *MockTaskRepository) UpdateTimeEntry(ctx context.Context, taskID string, timeEntryID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error) {
	args := m.Called(ctx, taskID, timeEntryID, timeEntry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TaskEntity), args.Error(1)
}

// BulkUpdate updates multiple tasks, handling conflicts when existing tasks are more recent
func (m *MockTaskRepository) BulkUpdate(ctx context.Context, tasks []*models.TaskEntity) ([]*models.TaskEntity, []*models.ConflictedItem, error) {
	args := m.Called(ctx, tasks)
	var updated []*models.TaskEntity
	var conflicts []*models.ConflictedItem
	
	if args.Get(0) != nil {
		updated = args.Get(0).([]*models.TaskEntity)
	}
	if args.Get(1) != nil {
		conflicts = args.Get(1).([]*models.ConflictedItem)
	}
	
	return updated, conflicts, args.Error(2)
}
