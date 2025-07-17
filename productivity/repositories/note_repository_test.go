package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/tests/utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupNoteTest(t *testing.T) (NoteRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewNoteRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestNoteRepository(t *testing.T) {
	repo, cleanup := setupNoteTest(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("Create note", func(t *testing.T) {
		note := &models.NoteEntity{
			Title:   stringPtr("Test Note"),
			Content: stringPtr("This is a test note"),
			User:    primitive.NewObjectID(),
		}

		createdNote, err := repo.Create(ctx, note)
		assert.NoError(t, err)
		assert.NotNil(t, createdNote)
		assert.NotNil(t, createdNote.ID)
		assert.Equal(t, "Test Note", *createdNote.Title)
		assert.Equal(t, "This is a test note", *createdNote.Content)
		assert.True(t, time.Since(createdNote.CreatedAt.Time()) < time.Second)
		assert.True(t, time.Since(createdNote.UpdatedAt.Time()) < time.Second)
	})

	t.Run("GetByID note", func(t *testing.T) {
		// Create a note first
		note := &models.NoteEntity{
			Title:   stringPtr("Test Note"),
			Content: stringPtr("This is a test note"),
			User:    primitive.NewObjectID(),
		}

		createdNote, err := repo.Create(ctx, note)
		assert.NoError(t, err)

		// Get the note by ID
		fetchedNote, err := repo.GetByID(ctx, createdNote.ID.Hex())
		assert.NoError(t, err)
		assert.NotNil(t, fetchedNote)
		assert.Equal(t, createdNote.ID, fetchedNote.ID)
		assert.Equal(t, *createdNote.Title, *fetchedNote.Title)
		assert.Equal(t, *createdNote.Content, *fetchedNote.Content)
	})

	t.Run("GetByID note not found", func(t *testing.T) {
		fetchedNote, err := repo.GetByID(ctx, primitive.NewObjectID().Hex())
		assert.NoError(t, err)
		assert.Nil(t, fetchedNote)
	})

	t.Run("GetAll notes", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create multiple notes for the user
		note1 := &models.NoteEntity{
			Title:   stringPtr("Note 1"),
			Content: stringPtr("Content 1"),
			User:    userID,
		}
		note2 := &models.NoteEntity{
			Title:   stringPtr("Note 2"),
			Content: stringPtr("Content 2"),
			User:    userID,
		}

		_, err := repo.Create(ctx, note1)
		assert.NoError(t, err)
		_, err = repo.Create(ctx, note2)
		assert.NoError(t, err)

		// Get all notes for the user
		notes, err := repo.GetAll(ctx, &userID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(notes), 2)

		// Check that all notes belong to the user
		for _, note := range notes {
			assert.Equal(t, userID, note.User)
		}
	})

	t.Run("Update note", func(t *testing.T) {
		// Create a note first
		note := &models.NoteEntity{
			Title:   stringPtr("Original Title"),
			Content: stringPtr("Original Content"),
			User:    primitive.NewObjectID(),
		}

		createdNote, err := repo.Create(ctx, note)
		assert.NoError(t, err)

		// Add a small delay to ensure different timestamps
		time.Sleep(1 * time.Millisecond)

		// Update the note
		updatedNote := &models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
			User:    createdNote.User,
		}

		result, err := repo.Update(ctx, createdNote.ID.Hex(), updatedNote)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Title", *result.Title)
		assert.Equal(t, "Updated Content", *result.Content)
		assert.True(t, result.UpdatedAt.Time().After(result.CreatedAt.Time()) || result.UpdatedAt.Time().Equal(result.CreatedAt.Time()))
	})

	t.Run("Delete note", func(t *testing.T) {
		// Create a note first
		note := &models.NoteEntity{
			Title:   stringPtr("To Delete"),
			Content: stringPtr("This will be deleted"),
			User:    primitive.NewObjectID(),
		}

		createdNote, err := repo.Create(ctx, note)
		assert.NoError(t, err)

		// Delete the note
		err = repo.Delete(ctx, createdNote.ID.Hex())
		assert.NoError(t, err)

		// Verify it's deleted
		fetchedNote, err := repo.GetByID(ctx, createdNote.ID.Hex())
		assert.NoError(t, err)
		assert.Nil(t, fetchedNote)
	})

	t.Run("Delete note not found", func(t *testing.T) {
		err := repo.Delete(ctx, primitive.NewObjectID().Hex())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "note not found")
	})
}

func TestNoteRepository_DeleteByUserID(t *testing.T) {
	repo, cleanup := setupNoteTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create test users
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create notes for user 1
	note1 := &models.NoteEntity{
		Title:   stringPtr("Note 1 for User 1"),
		Content: stringPtr("Content 1"),
		User:    userID1,
	}
	createdNote1, err := repo.Create(ctx, note1)
	require.NoError(t, err)
	require.NotNil(t, createdNote1)

	note2 := &models.NoteEntity{
		Title:   stringPtr("Note 2 for User 1"),
		Content: stringPtr("Content 2"),
		User:    userID1,
	}
	createdNote2, err := repo.Create(ctx, note2)
	require.NoError(t, err)
	require.NotNil(t, createdNote2)

	// Create notes for user 2
	note3 := &models.NoteEntity{
		Title:   stringPtr("Note 1 for User 2"),
		Content: stringPtr("Content 3"),
		User:    userID2,
	}
	createdNote3, err := repo.Create(ctx, note3)
	require.NoError(t, err)
	require.NotNil(t, createdNote3)

	// Verify all notes exist
	allNotes, err := repo.GetAll(ctx, nil)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(allNotes), 3)

	// Count notes for each user before deletion
	user1NotesBefore, err := repo.GetAll(ctx, &userID1)
	require.NoError(t, err)
	user2NotesBefore, err := repo.GetAll(ctx, &userID2)
	require.NoError(t, err)
	assert.Len(t, user1NotesBefore, 2)
	assert.Len(t, user2NotesBefore, 1)

	// Delete all notes for user 1
	err = repo.DeleteByUserID(ctx, userID1)
	require.NoError(t, err)

	// Verify user 1's notes are gone but user 2's remain
	user1NotesAfter, err := repo.GetAll(ctx, &userID1)
	require.NoError(t, err)
	user2NotesAfter, err := repo.GetAll(ctx, &userID2)
	require.NoError(t, err)
	assert.Len(t, user1NotesAfter, 0)
	assert.Len(t, user2NotesAfter, 1)
	assert.Equal(t, "Note 1 for User 2", *user2NotesAfter[0].Title)
	assert.Equal(t, userID2, user2NotesAfter[0].User)

	// Delete all notes for user 2
	err = repo.DeleteByUserID(ctx, userID2)
	require.NoError(t, err)

	// Verify no notes remain for user 2
	finalUser2Notes, err := repo.GetAll(ctx, &userID2)
	require.NoError(t, err)
	assert.Len(t, finalUser2Notes, 0)
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
