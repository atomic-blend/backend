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
