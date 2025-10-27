package user

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestDB(t *testing.T) (*Repository, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Get MongoDB connection URI
	mongoURI := mongoServer.URI()

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	require.NoError(t, err)

	// Get database reference and create repository
	db := client.Database("test_db")
	repo := NewUserRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestUserRepository_Create(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	email := "test@example.com"
	user := &models.UserEntity{
		Email: &email,
	}

	// Test creation
	created, err := repo.Create(context.Background(), user)
	assert.NoError(t, err)
	assert.NotNil(t, created.ID)
	assert.NotNil(t, created.CreatedAt)
	assert.NotNil(t, created.UpdatedAt)
	assert.Equal(t, email, *created.Email)
}

func TestUserRepository_FindByID(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test user first
	email := "test@example.com"
	user := &models.UserEntity{
		Email: &email,
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test finding by ID
	ginCtx := &gin.Context{}
	found, err := repo.FindByID(ginCtx, *created.ID)
	assert.NoError(t, err)
	assert.Equal(t, *created.ID, *found.ID)
	assert.Equal(t, *created.Email, *found.Email)

	// Test with non-existent ID
	nonExistentID := primitive.NewObjectID()
	_, err = repo.FindByID(ginCtx, nonExistentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test user first
	email := "test@example.com"
	user := &models.UserEntity{
		Email: &email,
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test finding by email
	found, err := repo.FindByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, *created.ID, *found.ID)
	assert.Equal(t, email, *found.Email)

	// Test with non-existent email
	_, err = repo.FindByEmail(context.Background(), "nonexistent@example.com")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Test with invalid email format
	_, err = repo.FindByEmail(context.Background(), "invalid-email-format")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid email format")
}

func TestUserRepository_Update(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test user first
	email := "test@example.com"
	user := &models.UserEntity{
		Email: &email,
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Store the original timestamp
	originalUpdatedAt := created.UpdatedAt

	// Wait a moment to ensure timestamp will be different
	time.Sleep(time.Millisecond)

	// Update the user
	newEmail := "updated@example.com"
	created.Email = &newEmail
	updated, err := repo.Update(context.Background(), created)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, *updated.Email)
	assert.NotEqual(t, originalUpdatedAt.Time(), updated.UpdatedAt.Time())

	// Test update with non-existent ID
	nonExistentID := primitive.NewObjectID()
	created.ID = &nonExistentID
	_, err = repo.Update(context.Background(), created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a test user first
	email := "test@example.com"
	user := &models.UserEntity{
		Email: &email,
	}
	created, err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	// Test deletion
	err = repo.Delete(context.Background(), created.ID.Hex())
	assert.NoError(t, err)

	// Verify user is deleted
	ginCtx := &gin.Context{}
	_, err = repo.FindByID(ginCtx, *created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	// Test delete with non-existent ID
	err = repo.Delete(context.Background(), primitive.NewObjectID().Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_FindInactiveSubscriptionUsers(t *testing.T) {
	repo, cleanup := setupTestDB(t)
	defer cleanup()

	now := time.Now()
	cutoffDate := now.AddDate(0, 0, -7) // 7 days ago

	// Create users with different scenarios

	// User 1: No subscriptionId, created more than 7 days ago (should be found)
	email1 := "inactive1@example.com"
	user1 := &models.UserEntity{
		Email: &email1,
	}
	created1, err := repo.Create(context.Background(), user1)
	require.NoError(t, err)
	// Set createdAt to 8 days ago
	pastCreatedAt1 := primitive.NewDateTimeFromTime(cutoffDate.AddDate(0, 0, -1))
	created1.CreatedAt = &pastCreatedAt1
	_, err = repo.Update(context.Background(), created1)
	require.NoError(t, err)

	// User 2: No subscriptionId, created recently (should not be found)
	email2 := "active@example.com"
	user2 := &models.UserEntity{
		Email: &email2,
	}
	created2, err := repo.Create(context.Background(), user2)
	require.NoError(t, err)
	// CreatedAt is recent, no need to change

	// User 3: Cancelled subscription, cancelled more than 7 days ago (should be found)
	email3 := "cancelled_old@example.com"
	status3 := "cancelled"
	subID3 := "sub_cancelled"
	user3 := &models.UserEntity{
		Email: &email3,
	}
	created3, err := repo.Create(context.Background(), user3)
	require.NoError(t, err)
	// Set createdAt to 10 days ago
	pastCreatedAt3 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -10))
	created3.CreatedAt = &pastCreatedAt3
	created3.SubscriptionStatus = &status3
	created3.StripeSubscriptionID = &subID3
	pastCancelledAt3 := primitive.NewDateTimeFromTime(cutoffDate.AddDate(0, 0, -1))
	created3.CancelledAt = &pastCancelledAt3
	_, err = repo.Update(context.Background(), created3)
	require.NoError(t, err)

	// User 4: Cancelled subscription, cancelled recently (should not be found)
	email4 := "cancelled_recent@example.com"
	status4 := "cancelled"
	subID4 := "sub_cancelled_recent"
	user4 := &models.UserEntity{
		Email: &email4,
	}
	created4, err := repo.Create(context.Background(), user4)
	require.NoError(t, err)
	// Set createdAt to 10 days ago
	pastCreatedAt4 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -10))
	created4.CreatedAt = &pastCreatedAt4
	created4.SubscriptionStatus = &status4
	created4.StripeSubscriptionID = &subID4
	recentCancelledAt4 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -1))
	created4.CancelledAt = &recentCancelledAt4
	_, err = repo.Update(context.Background(), created4)
	require.NoError(t, err)

	// User 5: Active subscription (should not be found)
	email5 := "active_sub@example.com"
	subID5 := "sub_123"
	user5 := &models.UserEntity{
		Email:                &email5,
		StripeSubscriptionID: &subID5,
	}
	created5, err := repo.Create(context.Background(), user5)
	require.NoError(t, err)
	// Set createdAt to 8 days ago
	pastCreatedAt5 := primitive.NewDateTimeFromTime(cutoffDate.AddDate(0, 0, -1))
	created5.CreatedAt = &pastCreatedAt5
	_, err = repo.Update(context.Background(), created5)
	require.NoError(t, err)

	// Call FindInactiveSubscriptionUsers with gracePeriodDays = 7
	inactiveUsers, err := repo.FindInactiveSubscriptionUsers(context.Background(), 7)
	assert.NoError(t, err)

	// Should find user1 and user3
	assert.Len(t, inactiveUsers, 2)

	// Check that the found users are user1 and user3
	foundIDs := make(map[string]bool)
	for _, u := range inactiveUsers {
		foundIDs[u.ID.Hex()] = true
	}

	assert.True(t, foundIDs[created1.ID.Hex()], "User1 should be found")
	assert.True(t, foundIDs[created3.ID.Hex()], "User3 should be found")
	assert.False(t, foundIDs[created2.ID.Hex()], "User2 should not be found")
	assert.False(t, foundIDs[created4.ID.Hex()], "User4 should not be found")
	assert.False(t, foundIDs[created5.ID.Hex()], "User5 should not be found")
}
