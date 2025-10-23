package stripe

import (
	"errors"
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetOrCreateCustomer(t *testing.T) {
	// Setup
	mockUserRepo := &mocks.MockUserRepository{}
	mockStripeClient := &mocks.MockStripeClient{}
	service := &Service{
		userService:  mockUserRepo,
		stripeClient: mockStripeClient,
	}
	ctx := &gin.Context{}

	t.Run("user not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("create new customer successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		user := &models.UserEntity{
			ID:        &userID,
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
		}
		customerID := "cus_123"
		customer := &stripe.Customer{
			ID: customerID,
		}
		updatedUser := &models.UserEntity{
			ID:               &userID,
			StripeCustomerId: &customerID,
			Email:            &email,
			FirstName:        &firstName,
			LastName:         &lastName,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockStripeClient.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*stripe.CustomerCreateParams")).Return(customer, nil).Once()
		mockUserRepo.On("Update", ctx, mock.AnythingOfType("*models.UserEntity")).Return(updatedUser, nil).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.NotNil(t, result)
		assert.Equal(t, customerID, result.ID)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("create customer error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		user := &models.UserEntity{
			ID:        &userID,
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockStripeClient.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*stripe.CustomerCreateParams")).Return(nil, errors.New("stripe error")).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("update user error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		user := &models.UserEntity{
			ID:        &userID,
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
		}
		customerID := "cus_123"
		customer := &stripe.Customer{
			ID: customerID,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockStripeClient.On("CreateCustomer", mock.Anything, mock.AnythingOfType("*stripe.CustomerCreateParams")).Return(customer, nil).Once()
		mockUserRepo.On("Update", ctx, mock.AnythingOfType("*models.UserEntity")).Return(nil, errors.New("update error")).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("get existing customer successfully", func(t *testing.T) {
		userID := primitive.NewObjectID()
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		customerID := "cus_123"
		user := &models.UserEntity{
			ID:               &userID,
			StripeCustomerId: &customerID,
			Email:            &email,
			FirstName:        &firstName,
			LastName:         &lastName,
		}
		customer := &stripe.Customer{
			ID: customerID,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockStripeClient.On("GetCustomer", mock.Anything, customerID).Return(customer, nil).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.NotNil(t, result)
		assert.Equal(t, customerID, result.ID)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("get customer error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"
		customerID := "cus_123"
		user := &models.UserEntity{
			ID:               &userID,
			StripeCustomerId: &customerID,
			Email:            &email,
			FirstName:        &firstName,
			LastName:         &lastName,
		}

		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockStripeClient.On("GetCustomer", mock.Anything, customerID).Return(nil, errors.New("get error")).Once()

		result := service.GetOrCreateCustomer(ctx, userID)

		assert.Nil(t, result)
		mockUserRepo.AssertExpectations(t)
		mockStripeClient.AssertExpectations(t)
	})
}

func TestCreateSubscription(t *testing.T) {
	// Setup
	mockUserRepo := &mocks.MockUserRepository{}
	mockStripeClient := &mocks.MockStripeClient{}
	service := &Service{
		userService:  mockUserRepo,
		stripeClient: mockStripeClient,
	}
	ctx := &gin.Context{}

	t.Run("successful subscription creation", func(t *testing.T) {
		customerID := "cus_123"
		priceID := "price_456"
		subscriptionID := "sub_789"
		subscription := &stripe.Subscription{
			ID: subscriptionID,
		}

		mockStripeClient.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*stripe.SubscriptionCreateParams")).Return(subscription, nil).Once()

		result := service.CreateSubscription(ctx, customerID, priceID)

		assert.NotNil(t, result)
		assert.Equal(t, subscriptionID, result.ID)
		mockStripeClient.AssertExpectations(t)
	})

	t.Run("subscription creation error", func(t *testing.T) {
		customerID := "cus_123"
		priceID := "price_456"

		mockStripeClient.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*stripe.SubscriptionCreateParams")).Return(nil, errors.New("stripe error")).Once()

		result := service.CreateSubscription(ctx, customerID, priceID)

		assert.Nil(t, result)
		mockStripeClient.AssertExpectations(t)
	})
}
