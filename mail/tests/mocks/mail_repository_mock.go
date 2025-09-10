package mocks

import (
	"context"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMailRepository provides a mock implementation of MailRepositoryInterface
type MockMailRepository struct {
	mock.Mock
}

// GetAll retrieves mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all mails and total count.
func (m *MockMailRepository) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.Mail, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.Mail), args.Get(1).(int64), args.Error(2)
}

// GetByID retrieves a mail by its ID
func (m *MockMailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Mail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Mail), args.Error(1)
}

// Create creates a new mail
func (m *MockMailRepository) Create(ctx context.Context, mail *models.Mail) (*models.Mail, error) {
	args := m.Called(ctx, mail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Mail), args.Error(1)
}

// CreateMany creates multiple mails
func (m *MockMailRepository) CreateMany(ctx context.Context, mails []models.Mail) (bool, error) {
	args := m.Called(ctx, mails)
	return args.Bool(0), args.Error(1)
}

// Update updates a mail object. This method is used by controller tests to mock mail updates.
func (m *MockMailRepository) Update(ctx context.Context, mail *models.Mail) error {
	args := m.Called(ctx, mail)
	return args.Error(0)
}
