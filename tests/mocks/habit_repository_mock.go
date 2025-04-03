package mocks

import (
	"atomic_blend_api/models"
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockHabitRepository provides a mock implementation of HabitRepositoryInterface
type MockHabitRepository struct {
	mock.Mock
}

// GetAll gets all habits
func (m *MockHabitRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Habit, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Habit), args.Error(1)
}

// GetByID gets a habit by ID
func (m *MockHabitRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Habit, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Habit), args.Error(1)
}

// Create creates a new habit
func (m *MockHabitRepository) Create(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	args := m.Called(ctx, habit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Habit), args.Error(1)
}

// Update updates a habit
func (m *MockHabitRepository) Update(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	args := m.Called(ctx, habit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Habit), args.Error(1)
}

// Delete deletes a habit by ID
func (m *MockHabitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
