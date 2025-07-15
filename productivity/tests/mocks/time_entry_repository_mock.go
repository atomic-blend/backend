package mocks

import (
	"productivity/models"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTimeEntryRepository is a mock implementation of TimeEntryRepositoryInterface
type MockTimeEntryRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockTimeEntryRepository) Create(ctx context.Context, timeEntry *models.TimeEntry) (*models.TimeEntry, error) {
	args := m.Called(ctx, timeEntry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeEntry), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockTimeEntryRepository) GetByID(ctx context.Context, id string) (*models.TimeEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeEntry), args.Error(1)
}

// GetAll mocks the GetAll method
func (m *MockTimeEntryRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TimeEntry, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TimeEntry), args.Error(1)
}

// Update mocks the Update method
func (m *MockTimeEntryRepository) Update(ctx context.Context, id string, timeEntry *models.TimeEntry) (*models.TimeEntry, error) {
	args := m.Called(ctx, id, timeEntry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TimeEntry), args.Error(1)
}

// Delete mocks the Delete method
func (m *MockTimeEntryRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
