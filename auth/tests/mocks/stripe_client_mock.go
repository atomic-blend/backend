package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v83"
)

// MockStripeClient provides a mock implementation of StripeClientInterface
type MockStripeClient struct {
	mock.Mock
}

// CreateCustomer mocks the CreateCustomer method
func (m *MockStripeClient) CreateCustomer(ctx context.Context, params *stripe.CustomerCreateParams) (*stripe.Customer, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

// GetCustomer mocks the GetCustomer method
func (m *MockStripeClient) GetCustomer(ctx context.Context, id string) (*stripe.Customer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Customer), args.Error(1)
}

// CreateSubscription mocks the CreateSubscription method
func (m *MockStripeClient) CreateSubscription(ctx context.Context, params *stripe.SubscriptionCreateParams) (*stripe.Subscription, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Subscription), args.Error(1)
}
