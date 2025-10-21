package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupResetPasswordTest(t *testing.T) (UserResetPasswordRequestRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewUserResetPasswordRequestRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestResetPasswordRequest() *models.UserResetPassword {
	userID := primitive.NewObjectID()
	resetCode := "RESET123456"
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.UserResetPassword{
		UserID:    &userID,
		ResetCode: resetCode,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestUserResetPasswordRequestRepository_Create(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("successful create reset password request", func(t *testing.T) {
		resetRequest := createTestResetPasswordRequest()

		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)
		assert.NotNil(t, created.UserID)
		assert.Equal(t, resetRequest.UserID, created.UserID)
		assert.Equal(t, resetRequest.ResetCode, created.ResetCode)
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})

	t.Run("create with different reset codes", func(t *testing.T) {
		resetRequest1 := createTestResetPasswordRequest()
		resetRequest1.ResetCode = "CODE1"

		resetRequest2 := createTestResetPasswordRequest()
		resetRequest2.ResetCode = "CODE2"

		created1, err := repo.Create(context.Background(), resetRequest1)
		require.NoError(t, err)

		created2, err := repo.Create(context.Background(), resetRequest2)
		require.NoError(t, err)

		assert.NotEqual(t, *created1.UserID, *created2.UserID)
		assert.NotEqual(t, created1.ResetCode, created2.ResetCode)
	})

	t.Run("create with same user ID but different reset codes", func(t *testing.T) {
		userID := primitive.NewObjectID()

		resetRequest1 := createTestResetPasswordRequest()
		resetRequest1.UserID = &userID
		resetRequest1.ResetCode = "FIRST_CODE"

		created1, err := repo.Create(context.Background(), resetRequest1)
		require.NoError(t, err)
		assert.Equal(t, userID, *created1.UserID)

		// Second request with same user ID should fail due to unique constraint
		resetRequest2 := createTestResetPasswordRequest()
		resetRequest2.UserID = &userID
		resetRequest2.ResetCode = "SECOND_CODE"

		_, err = repo.Create(context.Background(), resetRequest2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")
	})
}

func TestUserResetPasswordRequestRepository_FindByResetCode(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("successful find reset password request by reset code", func(t *testing.T) {
		resetRequest := createTestResetPasswordRequest()
		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)

		found, err := repo.FindByResetCode(context.Background(), created.ResetCode)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.UserID, *found.UserID)
		assert.Equal(t, created.ResetCode, found.ResetCode)
		assert.Equal(t, created.CreatedAt, found.CreatedAt)
		assert.Equal(t, created.UpdatedAt, found.UpdatedAt)
	})

	t.Run("reset password request not found by reset code", func(t *testing.T) {
		found, err := repo.FindByResetCode(context.Background(), "NONEXISTENT_CODE")
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("find with empty reset code", func(t *testing.T) {
		found, err := repo.FindByResetCode(context.Background(), "")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestUserResetPasswordRequestRepository_FindByUserID(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("successful find reset password request by user ID", func(t *testing.T) {
		resetRequest := createTestResetPasswordRequest()
		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)

		found, err := repo.FindByUserID(context.Background(), created.UserID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.UserID, *found.UserID)
		assert.Equal(t, created.ResetCode, found.ResetCode)
		assert.Equal(t, created.CreatedAt, found.CreatedAt)
		assert.Equal(t, created.UpdatedAt, found.UpdatedAt)
	})

	t.Run("reset password request not found by user ID", func(t *testing.T) {
		found, err := repo.FindByUserID(context.Background(), primitive.NewObjectID().Hex())
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("find with invalid user ID", func(t *testing.T) {
		found, err := repo.FindByUserID(context.Background(), "invalid-id")
		require.Error(t, err)
		assert.Nil(t, found)
	})

	t.Run("find with empty user ID", func(t *testing.T) {
		found, err := repo.FindByUserID(context.Background(), "")
		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestUserResetPasswordRequestRepository_Delete(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("successful delete reset password request by ID", func(t *testing.T) {
		resetRequest := createTestResetPasswordRequest()
		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), created.UserID.Hex())
		require.NoError(t, err)

		// Verify it's deleted by trying to find it
		found, err := repo.FindByUserID(context.Background(), created.UserID.Hex())
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("delete non-existent reset password request", func(t *testing.T) {
		err := repo.Delete(context.Background(), primitive.NewObjectID().Hex())
		require.NoError(t, err)
	})

	t.Run("delete with invalid ID", func(t *testing.T) {
		err := repo.Delete(context.Background(), "invalid-id")
		require.Error(t, err)
	})

	t.Run("delete with empty ID", func(t *testing.T) {
		err := repo.Delete(context.Background(), "")
		require.Error(t, err)
	})
}

func TestUserResetPasswordRequestRepository_Integration(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("complete CRUD operations", func(t *testing.T) {
		// Create
		resetRequest := createTestResetPasswordRequest()
		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)
		assert.NotNil(t, created.UserID)

		// Read by reset code
		foundByCode, err := repo.FindByResetCode(context.Background(), created.ResetCode)
		require.NoError(t, err)
		assert.NotNil(t, foundByCode)
		assert.Equal(t, *created.UserID, *foundByCode.UserID)

		// Read by user ID
		foundByUserID, err := repo.FindByUserID(context.Background(), created.UserID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, foundByUserID)
		assert.Equal(t, created.ResetCode, foundByUserID.ResetCode)

		// Delete
		err = repo.Delete(context.Background(), created.UserID.Hex())
		require.NoError(t, err)

		// Verify deletion
		deletedByCode, err := repo.FindByResetCode(context.Background(), created.ResetCode)
		require.NoError(t, err)
		assert.Nil(t, deletedByCode)

		deletedByUserID, err := repo.FindByUserID(context.Background(), created.UserID.Hex())
		require.NoError(t, err)
		assert.Nil(t, deletedByUserID)
	})

	t.Run("multiple reset requests for same user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create first reset request
		resetRequest1 := createTestResetPasswordRequest()
		resetRequest1.UserID = &userID
		resetRequest1.ResetCode = "FIRST_RESET_CODE"

		created1, err := repo.Create(context.Background(), resetRequest1)
		require.NoError(t, err)

		// Create second reset request for same user should fail due to unique constraint
		resetRequest2 := createTestResetPasswordRequest()
		resetRequest2.UserID = &userID
		resetRequest2.ResetCode = "SECOND_RESET_CODE"

		_, err = repo.Create(context.Background(), resetRequest2)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key")

		// Only the first request should exist
		found1, err := repo.FindByResetCode(context.Background(), created1.ResetCode)
		require.NoError(t, err)
		assert.NotNil(t, found1)
		assert.Equal(t, created1.ResetCode, found1.ResetCode)

		// Second reset code should not be found
		found2, err := repo.FindByResetCode(context.Background(), "SECOND_RESET_CODE")
		require.NoError(t, err)
		assert.Nil(t, found2)

		// Find by user ID should return the first one
		foundByUserID, err := repo.FindByUserID(context.Background(), userID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, foundByUserID)
		assert.Equal(t, userID, *foundByUserID.UserID)
		assert.Equal(t, created1.ResetCode, foundByUserID.ResetCode)
	})

	t.Run("reset code uniqueness", func(t *testing.T) {
		// Create first reset request
		resetRequest1 := createTestResetPasswordRequest()
		resetRequest1.ResetCode = "UNIQUE_CODE_123"

		created1, err := repo.Create(context.Background(), resetRequest1)
		require.NoError(t, err)

		// Create second reset request with different user but same reset code
		resetRequest2 := createTestResetPasswordRequest()
		resetRequest2.ResetCode = "UNIQUE_CODE_123"

		created2, err := repo.Create(context.Background(), resetRequest2)
		require.NoError(t, err)

		// Both should exist (no unique constraint on reset code in this implementation)
		assert.NotEqual(t, *created1.UserID, *created2.UserID)
		assert.Equal(t, created1.ResetCode, created2.ResetCode)

		// Find by reset code should return the first one found
		found, err := repo.FindByResetCode(context.Background(), "UNIQUE_CODE_123")
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "UNIQUE_CODE_123", found.ResetCode)
	})
}

func TestUserResetPasswordRequestRepository_EdgeCases(t *testing.T) {
	repo, cleanup := setupResetPasswordTest(t)
	defer cleanup()

	t.Run("very long reset code", func(t *testing.T) {
		longResetCode := "VERY_LONG_RESET_CODE_" + string(make([]byte, 1000))
		resetRequest := createTestResetPasswordRequest()
		resetRequest.ResetCode = longResetCode

		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)
		assert.Equal(t, longResetCode, created.ResetCode)

		found, err := repo.FindByResetCode(context.Background(), longResetCode)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, longResetCode, found.ResetCode)
	})

	t.Run("special characters in reset code", func(t *testing.T) {
		specialResetCode := "RESET-CODE_123!@#$%^&*()"
		resetRequest := createTestResetPasswordRequest()
		resetRequest.ResetCode = specialResetCode

		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)
		assert.Equal(t, specialResetCode, created.ResetCode)

		found, err := repo.FindByResetCode(context.Background(), specialResetCode)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, specialResetCode, found.ResetCode)
	})

	t.Run("unicode characters in reset code", func(t *testing.T) {
		unicodeResetCode := "RESET_CODE_🚀_测试_αβγ"
		resetRequest := createTestResetPasswordRequest()
		resetRequest.ResetCode = unicodeResetCode

		created, err := repo.Create(context.Background(), resetRequest)
		require.NoError(t, err)
		assert.Equal(t, unicodeResetCode, created.ResetCode)

		found, err := repo.FindByResetCode(context.Background(), unicodeResetCode)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, unicodeResetCode, found.ResetCode)
	})
}
