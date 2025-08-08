package mocks

import (
	"context"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockSendMailRepository provides a mock implementation of SendMailRepositoryInterface
type MockSendMailRepository struct {
	mock.Mock
}

// GetAll retrieves send mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all send mails and total count.
func (m *MockSendMailRepository) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.SendMail, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.SendMail), args.Get(1).(int64), args.Error(2)
}

// GetByID retrieves a send mail by its ID
func (m *MockSendMailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.SendMail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// Create creates a new send mail
func (m *MockSendMailRepository) Create(ctx context.Context, sendMail *models.SendMail) (*models.SendMail, error) {
	args := m.Called(ctx, sendMail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// UpdateStatus updates the status of a send mail
func (m *MockSendMailRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.SendStatus) (*models.SendMail, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// IncrementRetryCounter increments the retry counter of a send mail and sets status to retry
func (m *MockSendMailRepository) IncrementRetryCounter(ctx context.Context, id primitive.ObjectID) (*models.SendMail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// Delete soft deletes a send mail by marking it as trashed
func (m *MockSendMailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
