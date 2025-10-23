package payment

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name: "Successful subscription",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				customer := &stripe.Customer{ID: "cus_123"}
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{}, // Empty body on success
		},
		{
			name: "Stripe customer creation/retrieval failed",
			setupAuth: func(c *gin.Context) {
				userID := primitive.NewObjectID()
				c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "cannot_get_stripe_customer"},
		},
		{
			name: "Unauthorized access - no auth user",
			setupAuth: func(c *gin.Context) {
				// No auth user set
			},
			setupMocks: func(stripeService *mocks.MockStripeService) {
				// No mocks needed
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]string{"error": "Authentication required"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
			var responseBody map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &responseBody)
			if len(tc.expectedBody) == 0 {
				// For success case, body should be empty
				assert.Empty(t, w.Body.String())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, responseBody)
			}

			// Assert mocks
			mockStripeService.AssertExpectations(t)
		})
	}
}
