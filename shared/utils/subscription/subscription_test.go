package subscription

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestDB(t *testing.T) (*user.Repository, *gin.Context, func()) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	require.NoError(t, err)

	// Get database reference and create repository
	database := client.Database("test_db")
	repo := user.NewUserRepository(database)

	// Set the global database for the subscription function to use
	db.Database = database

	// Create gin context
	ctx := &gin.Context{}

	// Return cleanup function
	cleanup := func() {
		// Reset global database
		db.Database = nil
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, ctx, cleanup
}

func createTestUser(t *testing.T, repo *user.Repository, purchases []*models.PurchaseEntity) *models.UserEntity {
	email := "test@example.com"
	password := "testpassword"
	user := &models.UserEntity{
		Email:     &email,
		Password:  &password,
		Purchases: purchases,
	}

	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)
	require.NotNil(t, created.ID)

	return created
}

func TestIsUserSubscribed_UserNotFound(t *testing.T) {
	_, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Test with non-existent user ID
	nonExistentID := primitive.NewObjectID()
	result := IsUserSubscribed(ctx, nonExistentID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithoutPurchases(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create user without purchases
	user := createTestUser(t, repo, nil)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithExpiredSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with expired subscription
	expiredTime := time.Now().Add(-24 * time.Hour).UnixMilli() // 1 day ago
	purchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: expiredTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	purchase := models.NewRevenueCatPurchase(purchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithActiveSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with active subscription
	futureTime := time.Now().Add(24 * time.Hour).UnixMilli() // 1 day in future
	purchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: futureTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	purchase := models.NewRevenueCatPurchase(purchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithMultiplePurchases_OneActive(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple purchases - one expired, one active
	expiredTime := time.Now().Add(-24 * time.Hour).UnixMilli() // 1 day ago
	futureTime := time.Now().Add(24 * time.Hour).UnixMilli()   // 1 day in future

	expiredPurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: expiredTime,
		ProductID:      "expired_product",
		AppUserID:      "test_user",
	}
	expiredPurchase := models.NewRevenueCatPurchase(expiredPurchaseData)

	activePurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: futureTime,
		ProductID:      "active_product",
		AppUserID:      "test_user",
	}
	activePurchase := models.NewRevenueCatPurchase(activePurchaseData)

	purchases := []*models.PurchaseEntity{&expiredPurchase, &activePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithZeroExpirationTime(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with zero expiration time (should be treated as expired)
	purchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: 0,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	purchase := models.NewRevenueCatPurchase(purchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithNegativeExpirationTime(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with negative expiration time (should be treated as expired)
	purchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: -1000,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	purchase := models.NewRevenueCatPurchase(purchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithActiveStripeSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with active Stripe subscription
	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "active",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	purchase := models.NewStripePurchase(stripePurchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithInactiveStripeSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with inactive Stripe subscription
	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "canceled",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	purchase := models.NewStripePurchase(stripePurchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithPastDueStripeSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with past_due Stripe subscription (should be treated as inactive)
	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "past_due",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	purchase := models.NewStripePurchase(stripePurchaseData)
	purchases := []*models.PurchaseEntity{&purchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithBothRevenueCatAndStripe_BothActive(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create both RevenueCat and Stripe purchases, both active
	futureTime := time.Now().Add(24 * time.Hour).UnixMilli()
	rcPurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: futureTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	rcPurchase := models.NewRevenueCatPurchase(rcPurchaseData)

	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "active",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	stripePurchase := models.NewStripePurchase(stripePurchaseData)

	purchases := []*models.PurchaseEntity{&rcPurchase, &stripePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithBothRevenueCatAndStripe_OnlyStripeActive(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create both purchases - RevenueCat expired, Stripe active
	expiredTime := time.Now().Add(-24 * time.Hour).UnixMilli()
	rcPurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: expiredTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	rcPurchase := models.NewRevenueCatPurchase(rcPurchaseData)

	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "active",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	stripePurchase := models.NewStripePurchase(stripePurchaseData)

	purchases := []*models.PurchaseEntity{&rcPurchase, &stripePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithBothRevenueCatAndStripe_OnlyRevenueCatActive(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create both purchases - RevenueCat active, Stripe inactive
	futureTime := time.Now().Add(24 * time.Hour).UnixMilli()
	rcPurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: futureTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	rcPurchase := models.NewRevenueCatPurchase(rcPurchaseData)

	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "canceled",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	stripePurchase := models.NewStripePurchase(stripePurchaseData)

	purchases := []*models.PurchaseEntity{&rcPurchase, &stripePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithBothRevenueCatAndStripe_BothInactive(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create both purchases - both inactive
	expiredTime := time.Now().Add(-24 * time.Hour).UnixMilli()
	rcPurchaseData := models.RevenueCatPurchaseData{
		ExpirationAtMs: expiredTime,
		ProductID:      "test_product",
		AppUserID:      "test_user",
	}
	rcPurchase := models.NewRevenueCatPurchase(rcPurchaseData)

	stripePurchaseData := models.StripePurchaseData{
		ID:             "sub_123456",
		Status:         "incomplete_expired",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_123456",
		Amount:         999,
		Currency:       "usd",
	}
	stripePurchase := models.NewStripePurchase(stripePurchaseData)

	purchases := []*models.PurchaseEntity{&rcPurchase, &stripePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithMultipleStripePurchases(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create multiple Stripe purchases - one inactive, one active
	inactiveStripePurchaseData := models.StripePurchaseData{
		ID:             "sub_111111",
		Status:         "canceled",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_111111",
		Amount:         999,
		Currency:       "usd",
	}
	inactivePurchase := models.NewStripePurchase(inactiveStripePurchaseData)

	activeStripePurchaseData := models.StripePurchaseData{
		ID:             "sub_222222",
		Status:         "active",
		CustomerID:     "cus_123456",
		SubscriptionID: "sub_222222",
		Amount:         1999,
		Currency:       "usd",
	}
	activePurchase := models.NewStripePurchase(activeStripePurchaseData)

	purchases := []*models.PurchaseEntity{&inactivePurchase, &activePurchase}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}
