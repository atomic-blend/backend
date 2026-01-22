package payment

import (
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/auth/utils/stripe"

	"github.com/stretchr/testify/assert"
)

func TestNewController(t *testing.T) {
	// Create mock stripe service
	mockStripeService := new(mocks.MockStripeService)
	// Create mock user repository
	mockUserService := new(mocks.MockUserRepository)

	// Create controller
	controller := NewController(mockStripeService, mockUserService)

	// Assert controller properties
	assert.NotNil(t, controller)
	assert.Equal(t, mockStripeService, controller.stripeService)
}

func TestControllerImplementsInterfaces(t *testing.T) {
	// This is just a compile-time check to ensure the interfaces match
	var _ stripe.Interface = &mocks.MockStripeService{}
}
