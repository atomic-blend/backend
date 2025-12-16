package mocks

import (
	"context"
	"time"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCalendarRepository is a testify mock for calendar repository.
// It also contains a lightweight in-memory store used as a fallback when
// tests don't set explicit expectations (so existing tests keep working).
type MockCalendarRepository struct {
	mock.Mock
}

// GetAll retrieves all calendars for a user with pagination
func (m *MockCalendarRepository) GetAll(ctx *gin.Context, userID primitive.ObjectID, page int64, size int64) ([]*models.Calendar, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Calendar), args.Get(1).(int64), args.Error(2)
	}

	return args.Get(0).([]*models.Calendar), args.Get(1).(int64), args.Error(2)
}

// GetSince retrieves calendars for a user updated since a specific time with pagination
func (m *MockCalendarRepository) GetSince(ctx *gin.Context, userID primitive.ObjectID, since time.Time, page int64, size int64) ([]*models.Calendar, int64, error) {
	args := m.Called(ctx, userID, since, page, size)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Calendar), args.Get(1).(int64), args.Error(2)
	}

	return args.Get(0).([]*models.Calendar), args.Get(1).(int64), args.Error(2)
}

// GetByID retrieves a calendar by its ID
func (m *MockCalendarRepository) GetByID(ctx *gin.Context, id primitive.ObjectID) (*models.Calendar, error) {
	args := m.Called(ctx, id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}
	return nil, nil
}

// Create adds a new calendar
func (m *MockCalendarRepository) Create(ctx *gin.Context, calendar *models.Calendar) (*models.Calendar, error) {
	args := m.Called(ctx, calendar)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}

	id := primitive.NewObjectID()
	calendar.ID = &id
	now := primitive.NewDateTimeFromTime(time.Now())
	calendar.CreatedAt = &now
	calendar.UpdatedAt = &now
	return calendar, nil
}

// Update modifies an existing calendar
func (m *MockCalendarRepository) Update(ctx *gin.Context, id primitive.ObjectID, calendar *models.Calendar) (*models.Calendar, error) {
	args := m.Called(ctx, id, calendar)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}
	return args.Get(0).(*models.Calendar), args.Error(1)
}

// Delete removes a calendar by its ID
func (m *MockCalendarRepository) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	return nil
}

// CreateWithContext adds a new calendar using context.Context
func (m *MockCalendarRepository) CreateWithContext(ctx context.Context, calendar *models.Calendar) (*models.Calendar, error) {
	args := m.Called(ctx, calendar)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetByName retrieves a calendar by its name for a specific user
func (m *MockCalendarRepository) GetByName(ctx *gin.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetByNameWithContext retrieves a calendar by its name for a specific user using context.Context
func (m *MockCalendarRepository) GetByNameWithContext(ctx context.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Calendar), args.Error(1)
	}
	return nil, args.Error(1)
}