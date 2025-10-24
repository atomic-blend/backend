package mocks

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockStripeService provides a mock implementation of stripe.Interface
type MockStripeService struct {
	mock.Mock
}

// GetOrCreateCustomer mocks the GetOrCreateCustomer method
func (m *MockStripeService) GetOrCreateCustomer(ctx *gin.Context, userID primitive.ObjectID) *stripe.Customer {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.Customer)
}

// CreateSubscription mocks the CreateSubscription method
func (m *MockStripeService) CreateSubscription(ctx *gin.Context, customerID string, priceID string, trialDays int64) *stripe.Subscription {
	args := m.Called(ctx, customerID, priceID, trialDays)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.Subscription)
}

// GetSubscription mocks the GetSubscription method
func (m *MockStripeService) GetSubscription(ctx *gin.Context, customerID string, priceID string) *stripe.Subscription {
	args := m.Called(ctx, customerID, priceID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.Subscription)
}

// CreateInvoice mocks the CreateInvoice method
func (m *MockStripeService) CreateInvoice(ctx *gin.Context, customerID string, subscriptionID string) *stripe.Invoice {
	args := m.Called(ctx, customerID, subscriptionID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.Invoice)
}

// CreateInvoiceItem mocks the CreateInvoiceItem method
func (m *MockStripeService) CreateInvoiceItem(ctx *gin.Context, customerID string, amount float64, description string) *stripe.InvoiceItem {
	args := m.Called(ctx, customerID, amount, description)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.InvoiceItem)
}

// FinalizeInvoice mocks the FinalizeInvoice method
func (m *MockStripeService) FinalizeInvoice(ctx *gin.Context, invoiceID string) *stripe.Invoice {
	args := m.Called(ctx, invoiceID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.Invoice)
}

// GetEphemeralKeys mocks the GetEphemeralKeys method
func (m *MockStripeService) GetEphemeralKeys(ctx *gin.Context, customerID string) *stripe.EphemeralKey {
	args := m.Called(ctx, customerID)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*stripe.EphemeralKey)
}