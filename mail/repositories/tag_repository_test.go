package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTagTest(t *testing.T) (TagRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewTagRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestTag() *models.Tag {
	name := "Test Tag"
	color := "#FF5733"
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Tag{
		Name:      name,
		Color:     &color,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func TestTagRepository_Create(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	t.Run("successful create tag", func(t *testing.T) {
		tag := createTestTag()

		created, err := repo.Create(context.Background(), tag)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)
		assert.Equal(t, tag.Name, created.Name)
		assert.Equal(t, *tag.Color, *created.Color)
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})
}

func TestTagRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	t.Run("successful get tag", func(t *testing.T) {
		tag := createTestTag()
		created, err := repo.Create(context.Background(), tag)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)
		assert.Equal(t, tag.Name, found.Name)
		assert.Equal(t, *tag.Color, *found.Color)
	})

	t.Run("tag not found", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), primitive.NewObjectID())
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestTagRepository_GetAll(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	t.Run("successful get all tags for a user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create test tags for the user
		tag1 := createTestTag()
		tag1.UserID = &userID
		tag1.Name = "Tag 1"

		tag2 := createTestTag()
		tag2.UserID = &userID
		tag2.Name = "Tag 2"

		// Create one tag for another user
		otherUserID := primitive.NewObjectID()
		otherTag := createTestTag()
		otherTag.UserID = &otherUserID

		_, err := repo.Create(context.Background(), tag1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), tag2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), otherTag)
		require.NoError(t, err)

		// Get tags for the user
		tags, err := repo.GetAll(context.Background(), &userID)
		require.NoError(t, err)
		assert.Len(t, tags, 2)

		// Verify the tag names
		var names []string
		for _, t := range tags {
			names = append(names, t.Name)
		}
		assert.Contains(t, names, "Tag 1")
		assert.Contains(t, names, "Tag 2")

		// Get all tags (no user filter)
		allTags, err := repo.GetAll(context.Background(), nil)
		require.NoError(t, err)
		assert.Len(t, allTags, 3)
	})
}

func TestTagRepository_Update(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	t.Run("successful update tag", func(t *testing.T) {
		tag := createTestTag()
		created, err := repo.Create(context.Background(), tag)
		require.NoError(t, err)

		updatedName := "Updated Tag"
		updatedColor := "#00FF00"
		created.Name = updatedName
		created.Color = &updatedColor

		// Wait a moment to ensure timestamp will be different
		time.Sleep(time.Millisecond)

		updated, err := repo.Update(context.Background(), created)
		require.NoError(t, err)
		assert.Equal(t, updatedName, updated.Name)
		assert.Equal(t, updatedColor, *updated.Color)

		// Verify update is persisted
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedName, found.Name)
		assert.Equal(t, updatedColor, *found.Color)
	})

	t.Run("tag not found", func(t *testing.T) {
		tag := createTestTag()
		id := primitive.NewObjectID()
		tag.ID = &id

		_, err := repo.Update(context.Background(), tag)
		require.Error(t, err)
	})
}

func TestTagRepository_Delete(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	t.Run("successful delete tag", func(t *testing.T) {
		tag := createTestTag()
		created, err := repo.Create(context.Background(), tag)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), *created.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestTagRepository_DeleteByUserID(t *testing.T) {
	repo, cleanup := setupTagTest(t)
	defer cleanup()

	// Create test users
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create tags for user 1
	tag1 := createTestTag()
	tag1.UserID = &userID1
	tag1.Name = "Tag 1 for User 1"
	createdTag1, err := repo.Create(context.Background(), tag1)
	require.NoError(t, err)
	require.NotNil(t, createdTag1)

	tag2 := createTestTag()
	tag2.UserID = &userID1
	tag2.Name = "Tag 2 for User 1"
	createdTag2, err := repo.Create(context.Background(), tag2)
	require.NoError(t, err)
	require.NotNil(t, createdTag2)

	// Create tags for user 2
	tag3 := createTestTag()
	tag3.UserID = &userID2
	tag3.Name = "Tag 1 for User 2"
	createdTag3, err := repo.Create(context.Background(), tag3)
	require.NoError(t, err)
	require.NotNil(t, createdTag3)

	// Verify all tags exist
	allTags, err := repo.GetAll(context.Background(), nil)
	require.NoError(t, err)
	assert.Len(t, allTags, 3)

	// Count tags for each user before deletion
	user1TagsBefore, err := repo.GetAll(context.Background(), &userID1)
	require.NoError(t, err)
	user2TagsBefore, err := repo.GetAll(context.Background(), &userID2)
	require.NoError(t, err)
	assert.Len(t, user1TagsBefore, 2)
	assert.Len(t, user2TagsBefore, 1)

	// Delete all tags for user 1
	err = repo.DeleteByUserID(context.Background(), userID1)
	require.NoError(t, err)

	// Verify user 1's tags are gone but user 2's remain
	user1TagsAfter, err := repo.GetAll(context.Background(), &userID1)
	require.NoError(t, err)
	user2TagsAfter, err := repo.GetAll(context.Background(), &userID2)
	require.NoError(t, err)
	assert.Len(t, user1TagsAfter, 0)
	assert.Len(t, user2TagsAfter, 1)
	assert.Equal(t, "Tag 1 for User 2", user2TagsAfter[0].Name)
	assert.Equal(t, userID2, *user2TagsAfter[0].UserID)

	// Delete all tags for user 2
	err = repo.DeleteByUserID(context.Background(), userID2)
	require.NoError(t, err)

	// Verify no tags remain for user 2
	finalUser2Tags, err := repo.GetAll(context.Background(), &userID2)
	require.NoError(t, err)
	assert.Len(t, finalUser2Tags, 0)
}
