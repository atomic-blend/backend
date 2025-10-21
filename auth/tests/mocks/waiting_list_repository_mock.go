package mocks

import (
	"context"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/stretchr/testify/mock"
)

// MockWaitingListRepository provides a mock implementation of WaitingListRepositoryInterface
type MockWaitingListRepository struct {
	mock.Mock
}

// Count returns the number of waiting list records
func (m *MockWaitingListRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// CountWithCode returns the number of waiting list records that have a code
func (m *MockWaitingListRepository) CountWithCode(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

// Create creates a new waiting list record
func (m *MockWaitingListRepository) Create(ctx context.Context, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error) {
	args := m.Called(ctx, waitingList)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*waitinglist.WaitingList), args.Error(1)
}

// GetAll gets all waiting list records
func (m *MockWaitingListRepository) GetAll(ctx context.Context) ([]*waitinglist.WaitingList, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*waitinglist.WaitingList), args.Error(1)
}

// GetByID gets a waiting list record by ID
func (m *MockWaitingListRepository) GetByID(ctx context.Context, id string) (*waitinglist.WaitingList, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*waitinglist.WaitingList), args.Error(1)
}

// GetByEmail gets a waiting list record by email
func (m *MockWaitingListRepository) GetByEmail(ctx context.Context, email string) (*waitinglist.WaitingList, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*waitinglist.WaitingList), args.Error(1)
}

// GetByCode gets a waiting list record by code
func (m *MockWaitingListRepository) GetByCode(ctx context.Context, code string) (*waitinglist.WaitingList, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*waitinglist.WaitingList), args.Error(1)
}

// GetPositionByEmail gets the position of a waiting list record by email
func (m *MockWaitingListRepository) GetPositionByEmail(ctx context.Context, email string) (int64, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(int64), args.Error(1)
}

// Update updates a waiting list record with the given ID
func (m *MockWaitingListRepository) Update(ctx context.Context, id string, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error) {
	args := m.Called(ctx, id, waitingList)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*waitinglist.WaitingList), args.Error(1)
}

// Delete deletes a waiting list record by ID
func (m *MockWaitingListRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DeleteByEmail deletes a waiting list record by email
func (m *MockWaitingListRepository) DeleteByEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

// DeleteByCode deletes a waiting list record by code
func (m *MockWaitingListRepository) DeleteByCode(ctx context.Context, code string) error {
	args := m.Called(ctx, code)
	return args.Error(0)
}

// GetOldest gets the oldest N waiting list records
func (m *MockWaitingListRepository) GetOldest(ctx context.Context, limit int64) ([]*waitinglist.WaitingList, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*waitinglist.WaitingList), args.Error(1)
}
