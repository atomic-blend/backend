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

// GetAll retrieves all mails for a user
func (m *MockMailRepository) GetAll(ctx context.Context, userID primitive.ObjectID) ([]*models.Mail, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Mail), args.Error(1)
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
