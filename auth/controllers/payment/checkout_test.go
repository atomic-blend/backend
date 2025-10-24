package payment

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/atomic-blend/backend/auth/tests/mocks"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v83"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCheckout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockStripeService, *mocks.MockUserRepository)
		setupEnv       func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful checkout",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{}},
				}
				checkoutSession := &stripe.CheckoutSession{ID: "cs_123", URL: "https://checkout.stripe.com/pay/cs_123"}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
				stripeService.On("CreateCheckoutSession", mock.Anything, "cus_123", mock.AnythingOfType("int64")).Return(checkoutSession, nil)
			},
			setupEnv: func() {
				os.Unsetenv("STRIPE_CLOUD_TRIAL_DAYS")
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"session": "https://checkout.stripe.com/pay/cs_123",
			},
		},
		{
			name: "Stripe customer creation/retrieval failed",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, nil)
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_TRIAL_DAYS") },
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "cannot_get_stripe_customer"},
		},
		{
			name: "Unauthorized access - no auth user",
			setupAuth: func(c *gin.Context) {
				// No auth user set
			},
			setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
				// No mocks needed
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_TRIAL_DAYS") },
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Authentication required"},
		},
		{
			name: "subscription already exists",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{{ID: "sub_existing"}}},
				}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_TRIAL_DAYS") },
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "subscription_already_exists"},
		},
		{
			name: "Create checkout session failed",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{}},
				}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
				stripeService.On("CreateCheckoutSession", mock.Anything, "cus_123", mock.AnythingOfType("int64")).Return(nil, assert.AnError)
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_TRIAL_DAYS") },
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "cannot_create_checkout_session"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup env
			tc.setupEnv()

			// Create mocks
			mockStripeService := new(mocks.MockStripeService)
			mockUserService := new(mocks.MockUserRepository)

			// Setup mocks
			tc.setupMocks(mockStripeService, mockUserService)

			// Create controller
			controller := NewController(mockStripeService, mockUserService)

			// Create request
			req, _ := http.NewRequest("POST", "/payment/checkout", nil)

			// Create response recorder
			w := httptest.NewRecorder()

			// Create gin context
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			// Setup auth
			tc.setupAuth(ctx)

			// Call controller
			controller.Checkout(ctx)

			// Assert status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert body (compare only expected fields so tests aren't brittle
			// against full Stripe CheckoutSession serialization)
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			for key, expectedVal := range tc.expectedBody {
				assert.Equal(t, expectedVal, responseBody[key])
			}

			// Assert mocks
			mockStripeService.AssertExpectations(t)
			mockUserService.AssertExpectations(t)
		})
	}
}
