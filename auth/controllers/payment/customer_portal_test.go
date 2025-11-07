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

func TestCustomerPortal(t *testing.T) {
    gin.SetMode(gin.TestMode)

    testCases := []struct {
        name           string
        setupAuth      func(*gin.Context)
        setupMocks     func(*mocks.MockStripeService, *mocks.MockUserRepository)
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name: "Successful portal",
            setupAuth: func(c *gin.Context) {
                userID := primitive.NewObjectID()
                c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
            },
            setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
                customer := &stripe.Customer{ID: "cus_123"}
                portalSession := &stripe.BillingPortalSession{URL: "https://billing.stripe.com/session/bs_123"}
                stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
                stripeService.On("CreateCustomerPortalSession", mock.Anything, "cus_123", (*string)(nil)).Return(portalSession, nil)
            },
            expectedStatus: http.StatusOK,
            expectedBody:   map[string]interface{}{"url": "https://billing.stripe.com/session/bs_123"},
        },
        {
            name: "Unauthorized access - no auth user",
            setupAuth: func(c *gin.Context) {
                // no auth user
            },
            setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
            },
            expectedStatus: http.StatusUnauthorized,
            expectedBody:   map[string]interface{}{"error": "Authentication required"},
        },
        {
            name: "Stripe customer retrieval failed",
            setupAuth: func(c *gin.Context) {
                userID := primitive.NewObjectID()
                c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
            },
            setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
                stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil, nil)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "cannot_get_stripe_customer"},
        },
        {
            name: "Create customer portal session failed",
            setupAuth: func(c *gin.Context) {
                userID := primitive.NewObjectID()
                c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
            },
            setupMocks: func(stripeService *mocks.MockStripeService, userService *mocks.MockUserRepository) {
                customer := &stripe.Customer{ID: "cus_123"}
                stripeService.On("GetOrCreateCustomer", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(customer, nil)
                stripeService.On("CreateCustomerPortalSession", mock.Anything, "cus_123", (*string)(nil)).Return(nil, assert.AnError)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody:   map[string]interface{}{"error": "cannot_create_customer_portal_session"},
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create mocks
            mockStripeService := new(mocks.MockStripeService)
            mockUserService := new(mocks.MockUserRepository)

            // Setup mocks
            tc.setupMocks(mockStripeService, mockUserService)

            // Create controller
            controller := NewController(mockStripeService, mockUserService)

            // Create request
            req, _ := http.NewRequest("GET", "/payment/customer-portal", nil)

            // Create response recorder
            w := httptest.NewRecorder()

            // Create gin context
            ctx, _ := gin.CreateTestContext(w)
            ctx.Request = req

            // Setup auth
            tc.setupAuth(ctx)

            // Call controller
            controller.CustomerPortal(ctx)

            // Assert status
            assert.Equal(t, tc.expectedStatus, w.Code)

            // Assert body
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
