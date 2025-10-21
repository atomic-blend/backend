package repositories

import (
	"context"
	"testing"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupWaitingListTest(t *testing.T) (WaitingListRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewWaitingListRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestWaitingList() *waitinglist.WaitingList {
	email := "test@example.com"
	code := "TEST123"
	now := primitive.NewDateTimeFromTime(time.Now())
	return &waitinglist.WaitingList{
		Email:     email,
		Code:      &code,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func TestWaitingListRepository_Create(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful create waiting list", func(t *testing.T) {
		waitingList := createTestWaitingList()

		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)
		assert.Equal(t, waitingList.Email, created.Email)
		assert.Equal(t, *waitingList.Code, *created.Code)
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})

	t.Run("create with existing ID", func(t *testing.T) {
		id := primitive.NewObjectID()
		waitingList := createTestWaitingList()
		waitingList.ID = &id

		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)
		assert.Equal(t, id, *created.ID)
		assert.Equal(t, waitingList.Email, created.Email)
	})

	t.Run("create without code", func(t *testing.T) {
		waitingList := &waitinglist.WaitingList{
			Email: "test2@example.com",
		}

		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)
		assert.Equal(t, waitingList.Email, created.Email)
		assert.Nil(t, created.Code)
	})
}

func TestWaitingListRepository_GetByID(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful get waiting list by ID", func(t *testing.T) {
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), created.ID.Hex())
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)
		assert.Equal(t, waitingList.Email, found.Email)
		assert.Equal(t, *waitingList.Code, *found.Code)
	})

	t.Run("waiting list not found", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), primitive.NewObjectID().Hex())
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("invalid ObjectID", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), "invalid-id")
		require.Error(t, err)
		assert.Nil(t, found)
	})
}

func TestWaitingListRepository_GetByEmail(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful get waiting list by email", func(t *testing.T) {
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		found, err := repo.GetByEmail(context.Background(), waitingList.Email)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)
		assert.Equal(t, waitingList.Email, found.Email)
		assert.Equal(t, *waitingList.Code, *found.Code)
	})

	t.Run("waiting list not found by email", func(t *testing.T) {
		found, err := repo.GetByEmail(context.Background(), "nonexistent@example.com")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestWaitingListRepository_GetByCode(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful get waiting list by code", func(t *testing.T) {
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		found, err := repo.GetByCode(context.Background(), *waitingList.Code)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)
		assert.Equal(t, waitingList.Email, found.Email)
		assert.Equal(t, *waitingList.Code, *found.Code)
	})

	t.Run("waiting list not found by code", func(t *testing.T) {
		found, err := repo.GetByCode(context.Background(), "NONEXISTENT")
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestWaitingListRepository_GetAll(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful get all waiting lists", func(t *testing.T) {
		// Create multiple waiting lists
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "test1@example.com"
		waitingList1.Code = stringPtr("CODE1")

		waitingList2 := createTestWaitingList()
		waitingList2.Email = "test2@example.com"
		waitingList2.Code = stringPtr("CODE2")

		waitingList3 := createTestWaitingList()
		waitingList3.Email = "test3@example.com"
		waitingList3.Code = stringPtr("CODE3")

		_, err := repo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), waitingList3)
		require.NoError(t, err)

		// Get all waiting lists
		allWaitingLists, err := repo.GetAll(context.Background())
		require.NoError(t, err)
		assert.Len(t, allWaitingLists, 3)

		// Verify the emails are present
		var emails []string
		for _, wl := range allWaitingLists {
			emails = append(emails, wl.Email)
		}
		assert.Contains(t, emails, "test1@example.com")
		assert.Contains(t, emails, "test2@example.com")
		assert.Contains(t, emails, "test3@example.com")
	})

	t.Run("get all when no waiting lists exist", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		allWaitingLists, err := freshRepo.GetAll(context.Background())
		require.NoError(t, err)
		assert.Len(t, allWaitingLists, 0)
	})
}

func TestWaitingListRepository_Count(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("count with no records", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		count, err := freshRepo.Count(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("count with multiple records", func(t *testing.T) {
		// Create multiple waiting lists
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "test1@example.com"

		waitingList2 := createTestWaitingList()
		waitingList2.Email = "test2@example.com"

		waitingList3 := createTestWaitingList()
		waitingList3.Email = "test3@example.com"

		_, err := repo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), waitingList3)
		require.NoError(t, err)

		count, err := repo.Count(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})
}

func TestWaitingListRepository_CountWithCode(t *testing.T) {
	t.Run("count with code when no records exist", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		count, err := freshRepo.CountWithCode(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("count with code when all records have codes", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		// Create multiple waiting lists with codes
		code1 := "CODE123"
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "test1@example.com"
		waitingList1.Code = &code1

		code2 := "CODE456"
		waitingList2 := createTestWaitingList()
		waitingList2.Email = "test2@example.com"
		waitingList2.Code = &code2

		code3 := "CODE789"
		waitingList3 := createTestWaitingList()
		waitingList3.Email = "test3@example.com"
		waitingList3.Code = &code3

		_, err := freshRepo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		_, err = freshRepo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		_, err = freshRepo.Create(context.Background(), waitingList3)
		require.NoError(t, err)

		count, err := freshRepo.CountWithCode(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(3), count)
	})

	t.Run("count with code when some records have codes", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		// Create waiting lists with and without codes
		code1 := "CODE123"
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "test4@example.com"
		waitingList1.Code = &code1

		waitingList2 := createTestWaitingList()
		waitingList2.Email = "test5@example.com"
		waitingList2.Code = nil // No code

		code3 := "CODE789"
		waitingList3 := createTestWaitingList()
		waitingList3.Email = "test6@example.com"
		waitingList3.Code = &code3

		_, err := freshRepo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		_, err = freshRepo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		_, err = freshRepo.Create(context.Background(), waitingList3)
		require.NoError(t, err)

		count, err := freshRepo.CountWithCode(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(2), count) // Only 2 have codes
	})

	t.Run("count with code when no records have codes", func(t *testing.T) {
		// Use a fresh repository instance for this test
		freshRepo, freshCleanup := setupWaitingListTest(t)
		defer freshCleanup()

		// Create waiting lists without codes
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "test7@example.com"
		waitingList1.Code = nil

		waitingList2 := createTestWaitingList()
		waitingList2.Email = "test8@example.com"
		waitingList2.Code = nil

		_, err := freshRepo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		_, err = freshRepo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		count, err := freshRepo.CountWithCode(context.Background())
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestWaitingListRepository_Update(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful update waiting list", func(t *testing.T) {
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		// Update the waiting list
		newCode := "UPDATED123"
		created.Code = &newCode
		created.Email = "updated@example.com"

		updated, err := repo.Update(context.Background(), created.ID.Hex(), created)
		require.NoError(t, err)
		assert.Equal(t, *created.ID, *updated.ID)
		assert.Equal(t, "updated@example.com", updated.Email)
		assert.Equal(t, "UPDATED123", *updated.Code)
		assert.NotNil(t, updated.UpdatedAt)
	})

	t.Run("update non-existent waiting list", func(t *testing.T) {
		waitingList := createTestWaitingList()
		id := primitive.NewObjectID()
		waitingList.ID = &id

		updated, err := repo.Update(context.Background(), waitingList.ID.Hex(), waitingList)
		require.NoError(t, err)
		assert.NotNil(t, updated)
	})

	t.Run("update with invalid ObjectID", func(t *testing.T) {
		waitingList := createTestWaitingList()

		updated, err := repo.Update(context.Background(), "invalid-id", waitingList)
		require.Error(t, err)
		assert.Nil(t, updated)
	})
}

func TestWaitingListRepository_Delete(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful delete waiting list by ID", func(t *testing.T) {
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), created.ID.Hex())
		require.NoError(t, err)

		// Verify it's deleted
		found, err := repo.GetByID(context.Background(), created.ID.Hex())
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("delete non-existent waiting list", func(t *testing.T) {
		err := repo.Delete(context.Background(), primitive.NewObjectID().Hex())
		require.NoError(t, err)
	})

	t.Run("delete with invalid ObjectID", func(t *testing.T) {
		err := repo.Delete(context.Background(), "invalid-id")
		require.Error(t, err)
	})
}

func TestWaitingListRepository_DeleteByEmail(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful delete waiting list by email", func(t *testing.T) {
		waitingList := createTestWaitingList()
		_, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		err = repo.DeleteByEmail(context.Background(), waitingList.Email)
		require.NoError(t, err)

		// Verify it's deleted
		found, err := repo.GetByEmail(context.Background(), waitingList.Email)
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("delete non-existent waiting list by email", func(t *testing.T) {
		err := repo.DeleteByEmail(context.Background(), "nonexistent@example.com")
		require.NoError(t, err)
	})
}

func TestWaitingListRepository_DeleteByCode(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("successful delete waiting list by code", func(t *testing.T) {
		waitingList := createTestWaitingList()
		_, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)

		err = repo.DeleteByCode(context.Background(), *waitingList.Code)
		require.NoError(t, err)

		// Verify it's deleted
		found, err := repo.GetByCode(context.Background(), *waitingList.Code)
		require.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("delete non-existent waiting list by code", func(t *testing.T) {
		err := repo.DeleteByCode(context.Background(), "NONEXISTENT")
		require.NoError(t, err)
	})
}

func TestWaitingListRepository_Integration(t *testing.T) {
	repo, cleanup := setupWaitingListTest(t)
	defer cleanup()

	t.Run("complete CRUD operations", func(t *testing.T) {
		// Create
		waitingList := createTestWaitingList()
		created, err := repo.Create(context.Background(), waitingList)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)

		// Read
		found, err := repo.GetByID(context.Background(), created.ID.Hex())
		require.NoError(t, err)
		assert.Equal(t, created.Email, found.Email)

		// Update
		newEmail := "updated@example.com"
		found.Email = newEmail
		updated, err := repo.Update(context.Background(), found.ID.Hex(), found)
		require.NoError(t, err)
		assert.Equal(t, newEmail, updated.Email)

		// Verify update
		verifyUpdated, err := repo.GetByEmail(context.Background(), newEmail)
		require.NoError(t, err)
		assert.Equal(t, newEmail, verifyUpdated.Email)

		// Delete
		err = repo.Delete(context.Background(), updated.ID.Hex())
		require.NoError(t, err)

		// Verify deletion
		deleted, err := repo.GetByID(context.Background(), updated.ID.Hex())
		require.NoError(t, err)
		assert.Nil(t, deleted)
	})

	t.Run("duplicate email handling", func(t *testing.T) {
		// Create first waiting list
		waitingList1 := createTestWaitingList()
		waitingList1.Email = "duplicate@example.com"
		created1, err := repo.Create(context.Background(), waitingList1)
		require.NoError(t, err)

		// Try to create second waiting list with same email
		waitingList2 := createTestWaitingList()
		waitingList2.Email = "duplicate@example.com"
		waitingList2.Code = stringPtr("DIFFERENT_CODE")

		created2, err := repo.Create(context.Background(), waitingList2)
		require.NoError(t, err)

		// Both should exist (no unique constraint in this implementation)
		assert.NotEqual(t, *created1.ID, *created2.ID)
		assert.Equal(t, created1.Email, created2.Email)

		// Get by email should return the first one found
		found, err := repo.GetByEmail(context.Background(), "duplicate@example.com")
		require.NoError(t, err)
		assert.NotNil(t, found)
	})
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
