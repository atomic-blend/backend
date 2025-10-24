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

func TestSubscribe(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		setupAuth      func(*gin.Context)
		setupMocks     func(*mocks.MockStripeService)
		setupEnv       func()
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful subscription",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{}},
				}
				subscription := &stripe.Subscription{
					ID:                 "sub_789",
					PendingSetupIntent: &stripe.SetupIntent{ClientSecret: "seti_123_secret"},
				}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
				stripeService.On("CreateSubscription", mock.Anything, "cus_123", "price_456").Return(subscription, nil)
			},
			setupEnv: func() {
				os.Setenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID", "price_456")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"subscription": map[string]interface{}{"intent": "", "secret": "seti_123_secret"}},
		},
		{
			name: "Stripe customer creation/retrieval failed",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, nil)
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID") },
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "cannot_get_stripe_customer"},
		},
		{
			name: "Unauthorized access - no auth user",
			setupAuth: func(c *gin.Context) {
				// No auth user set
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				// No mocks needed
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID") },
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]interface{}{"error": "Authentication required"},
		},
		{
			name: "subscription already exists",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{{ID: "sub_existing"}}},
				}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
				stripeService.On("GetSubscription", mock.Anything, "cus_123", "").Return(nil, nil)
			},
			setupEnv:       func() { os.Unsetenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID") },
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "subscription_already_exists"},
		},
		{
			name: "Subscription already exists with pending setup intent",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				customer := &stripe.Customer{
					ID:            "cus_123",
					Subscriptions: &stripe.SubscriptionList{Data: []*stripe.Subscription{{ID: "sub_existing"}}},
				}
				subscription := &stripe.Subscription{
					ID:                 "sub_existing",
					PendingSetupIntent: &stripe.SetupIntent{ID: "seti_123", ClientSecret: "seti_123_secret"},
				}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
				stripeService.On("GetSubscription", mock.Anything, "cus_123", "price_456").Return(subscription, nil)
			},
			setupEnv: func() {
				os.Setenv("STRIPE_CLOUD_SUBSCRIPTION_PRICE_ID", "price_456")
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]interface{}{"subscription": map[string]interface{}{"intent": "seti_123", "secret": "seti_123_secret"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup environment
			tc.setupEnv()

			// Create mock stripe service
			mockStripeService := new(mocks.MockStripeService)

			// Setup mocks
			tc.setupMocks(mockStripeService)

			// Create controller
			controller := NewController(mockStripeService)

			// Create a test HTTP request
			req, _ := http.NewRequest("POST", "/payment/subscribe", nil)

			// Create a response recorder
			w := httptest.NewRecorder()

			// Create a new Gin context
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			// Setup auth
			tc.setupAuth(ctx)

			// Call the controller method
			controller.Subscribe(ctx)

			// Assert status code
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Assert response body
			var responseBody map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBody, responseBody)

			// Assert mocks
			mockStripeService.AssertExpectations(t)
		})
	}
}
