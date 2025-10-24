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
func (m *MockStripeClient) GetCustomer(ctx context.Context, id string, params *stripe.CustomerRetrieveParams) (*stripe.Customer, error) {
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

// GetSubscription mocks the GetSubscription method
func (m *MockStripeClient) GetSubscription(ctx context.Context, customerID string, priceID string) (*stripe.Subscription, error) {
	args := m.Called(ctx, customerID, priceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Subscription), args.Error(1)
}

// CreateInvoice mocks the CreateInvoice method
func (m *MockStripeClient) CreateInvoice(ctx context.Context, params *stripe.InvoiceCreateParams) (*stripe.Invoice, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Invoice), args.Error(1)
}

// CreateInvoiceItem mocks the CreateInvoiceItem method
func (m *MockStripeClient) CreateInvoiceItem(ctx context.Context, params *stripe.InvoiceItemCreateParams) (*stripe.InvoiceItem, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.InvoiceItem), args.Error(1)
}

// FinalizeInvoice mocks the FinalizeInvoice method
func (m *MockStripeClient) FinalizeInvoice(ctx context.Context, id string, params *stripe.InvoiceFinalizeInvoiceParams) (*stripe.Invoice, error) {
	args := m.Called(ctx, id, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.Invoice), args.Error(1)
}

// GetEphemeralKeys mocks the GetEphemeralKeys method
func (m *MockStripeClient) GetEphemeralKeys(ctx context.Context, params *stripe.EphemeralKeyCreateParams) (*stripe.EphemeralKey, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*stripe.EphemeralKey), args.Error(1)
}