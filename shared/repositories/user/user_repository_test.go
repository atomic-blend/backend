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

func setupTestDB(t *testing.T) (*UserRepository, func()) {
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
