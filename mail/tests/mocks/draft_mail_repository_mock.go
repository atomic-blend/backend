// Package mocks is a package that contains the mocks for the microservice
package mocks

import (
	"context"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/stretchr/testify/mock"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockDraftMailRepository provides a mock implementation of DraftMailRepositoryInterface
type MockDraftMailRepository struct {
	mock.Mock
}

// GetAll retrieves draft mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
func (m *MockDraftMailRepository) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.SendMail, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.SendMail), args.Get(1).(int64), args.Error(2)
}

// GetByID retrieves a draft mail by its ID
func (m *MockDraftMailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.SendMail, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// Create creates a new draft mail
func (m *MockDraftMailRepository) Create(ctx context.Context, sendMail *models.SendMail) (*models.SendMail, error) {
	args := m.Called(ctx, sendMail)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// Update updates a draft mail by its ID
func (m *MockDraftMailRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.SendMail, error) {
	args := m.Called(ctx, id, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SendMail), args.Error(1)
}

// Delete soft deletes a draft mail by marking it as trashed
func (m *MockDraftMailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Trash soft deletes a draft mail by marking it as trashed
func (m *MockDraftMailRepository) Trash(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// GetSince retrieves draft mails where updated_at is after the specified time. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
func (m *MockDraftMailRepository) GetSince(ctx context.Context, since time.Time, page, limit int64) ([]*models.SendMail, int64, error) {
	args := m.Called(ctx, since, page, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*models.SendMail), args.Get(1).(int64), args.Error(2)
}
