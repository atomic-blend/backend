package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupFolderTest(t *testing.T) (FolderRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewFolderRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestFolder() *models.Folder {
	name := "Test Folder"
	color := "#FF5733"
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Folder{
		Name:      name,
		Color:     &color,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func TestFolderRepository_Create(t *testing.T) {
	repo, cleanup := setupFolderTest(t)
	defer cleanup()

	t.Run("successful create folder", func(t *testing.T) {
		userID := primitive.NewObjectID()
		folder := createTestFolder()
		folder.UserID = userID

		created, err := repo.Create(context.Background(), folder)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)
		assert.Equal(t, folder.Name, created.Name)
		assert.Equal(t, *folder.Color, *created.Color)
		assert.Equal(t, userID, created.UserID)
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})
}

func TestFolderRepository_GetAll(t *testing.T) {
	repo, cleanup := setupFolderTest(t)
	defer cleanup()

	t.Run("successful get all folders for a user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create test folders for the user
		folder1 := createTestFolder()
		folder1.UserID = userID
		folder1.Name = "Folder 1"

		folder2 := createTestFolder()
		folder2.UserID = userID
		folder2.Name = "Folder 2"

		// Create one folder for another user
		otherUserID := primitive.NewObjectID()
		otherFolder := createTestFolder()
		otherFolder.UserID = otherUserID

		_, err := repo.Create(context.Background(), folder1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), folder2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), otherFolder)
		require.NoError(t, err)

		// Get folders for the user
		folders, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, folders, 2)

		// Verify the folder names
		var names []string
		for _, f := range folders {
			names = append(names, f.Name)
		}
		assert.Contains(t, names, "Folder 1")
		assert.Contains(t, names, "Folder 2")
	})
}

func TestFolderRepository_Update(t *testing.T) {
	repo, cleanup := setupFolderTest(t)
	defer cleanup()

	t.Run("successful update folder", func(t *testing.T) {
		userID := primitive.NewObjectID()
		folder := createTestFolder()
		folder.UserID = userID
		created, err := repo.Create(context.Background(), folder)
		require.NoError(t, err)

		updatedName := "Updated Folder"
		updatedColor := "#00FF00"
		created.Name = updatedName
		created.Color = &updatedColor

		// Wait a moment to ensure timestamp will be different
		time.Sleep(time.Millisecond)

		updated, err := repo.Update(context.Background(), *created.ID, created)
		require.NoError(t, err)
		assert.Equal(t, updatedName, updated.Name)
		assert.Equal(t, updatedColor, *updated.Color)
	})
}

func TestFolderRepository_Delete(t *testing.T) {
	repo, cleanup := setupFolderTest(t)
	defer cleanup()

	t.Run("successful delete folder", func(t *testing.T) {
		userID := primitive.NewObjectID()
		folder := createTestFolder()
		folder.UserID = userID
		created, err := repo.Create(context.Background(), folder)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), *created.ID)
		require.NoError(t, err)

		// Verify deletion by checking if user has any folders
		folders, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, folders, 0)
	})
}

func TestFolderRepository_DeleteByUserID(t *testing.T) {
	repo, cleanup := setupFolderTest(t)
	defer cleanup()

	// Create test users
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create folders for user 1
	folder1 := createTestFolder()
	folder1.UserID = userID1
	folder1.Name = "Folder 1 for User 1"
	createdFolder1, err := repo.Create(context.Background(), folder1)
	require.NoError(t, err)
	require.NotNil(t, createdFolder1)

	folder2 := createTestFolder()
	folder2.UserID = userID1
	folder2.Name = "Folder 2 for User 1"
	createdFolder2, err := repo.Create(context.Background(), folder2)
	require.NoError(t, err)
	require.NotNil(t, createdFolder2)

	// Create folders for user 2
	folder3 := createTestFolder()
	folder3.UserID = userID2
	folder3.Name = "Folder 1 for User 2"
	createdFolder3, err := repo.Create(context.Background(), folder3)
	require.NoError(t, err)
	require.NotNil(t, createdFolder3)

	// Count folders for each user before deletion
	user1FoldersBefore, err := repo.GetAll(context.Background(), userID1)
	require.NoError(t, err)
	user2FoldersBefore, err := repo.GetAll(context.Background(), userID2)
	require.NoError(t, err)
	assert.Len(t, user1FoldersBefore, 2)
	assert.Len(t, user2FoldersBefore, 1)

	// Delete all folders for user 1
	err = repo.DeleteByUserID(context.Background(), userID1)
	require.NoError(t, err)

	// Verify user 1's folders are gone but user 2's remain
	user1FoldersAfter, err := repo.GetAll(context.Background(), userID1)
	require.NoError(t, err)
	user2FoldersAfter, err := repo.GetAll(context.Background(), userID2)
	require.NoError(t, err)
	assert.Len(t, user1FoldersAfter, 0)
	assert.Len(t, user2FoldersAfter, 1)
	assert.Equal(t, "Folder 1 for User 2", user2FoldersAfter[0].Name)
	assert.Equal(t, userID2, user2FoldersAfter[0].UserID)

	// Delete all folders for user 2
	err = repo.DeleteByUserID(context.Background(), userID2)
	require.NoError(t, err)

	// Verify no folders remain for user 2
	finalUser2Folders, err := repo.GetAll(context.Background(), userID2)
	require.NoError(t, err)
	assert.Len(t, finalUser2Folders, 0)
}
