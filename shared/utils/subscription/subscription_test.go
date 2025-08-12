package subscription

import (
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/repositories/user"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/atomic-blend/backend/shared/utils/db"
	"context"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestDB(t *testing.T) (*user.UserRepository, *gin.Context, func()) {
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

func createTestUser(t *testing.T, repo *user.UserRepository, purchases []*models.PurchaseEntity) *models.UserEntity {
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
	purchaseType := "REVENUE_CAT"
	purchases := []*models.PurchaseEntity{
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: expiredTime,
				ProductID:      "test_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithActiveSubscription(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with active subscription
	futureTime := time.Now().Add(24 * time.Hour).UnixMilli() // 1 day in future
	purchaseType := "REVENUE_CAT"
	purchases := []*models.PurchaseEntity{
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: futureTime,
				ProductID:      "test_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

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
	purchaseType := "REVENUE_CAT"
	purchases := []*models.PurchaseEntity{
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: expiredTime,
				ProductID:      "expired_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: futureTime,
				ProductID:      "active_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.True(t, result)
}

func TestIsUserSubscribed_UserWithZeroExpirationTime(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with zero expiration time (should be treated as expired)
	purchaseType := "REVENUE_CAT"
	purchases := []*models.PurchaseEntity{
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: 0,
				ProductID:      "test_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}

func TestIsUserSubscribed_UserWithNegativeExpirationTime(t *testing.T) {
	repo, ctx, cleanup := setupTestDB(t)
	defer cleanup()

	// Create purchase with negative expiration time (should be treated as expired)
	purchaseType := "REVENUE_CAT"
	purchases := []*models.PurchaseEntity{
		{
			ID:   primitive.NewObjectID(),
			Type: &purchaseType,
			PurchaseData: models.RevenueCatPurchaseData{
				ExpirationAtMs: -1000,
				ProductID:      "test_product",
				AppUserID:      "test_user",
			},
			CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
		},
	}

	user := createTestUser(t, repo, purchases)

	result := IsUserSubscribed(ctx, *user.ID)
	assert.False(t, result)
}
