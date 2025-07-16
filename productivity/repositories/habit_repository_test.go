package repositories

import (
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/tests/utils/inmemorymongo"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupHabitTest(t *testing.T) (HabitRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewHabitRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestHabit() *models.Habit {
	userID := primitive.NewObjectID()
	name := "Test Habit"
	frequency := models.FrequencyDaily
	now := primitive.NewDateTimeFromTime(time.Now())
	emoji := "💪"
	return &models.Habit{
		UserID:    userID,
		Name:      &name,
		Emoji:     &emoji,
		Frequency: &frequency,
		StartDate: &now,
	}
}

func createTestHabitEntry() *models.HabitEntry {
	habitID := primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.HabitEntry{
		HabitID:   habitID,
		EntryDate: now,
	}
}

func TestHabitRepository_Create(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful create habit", func(t *testing.T) {
		habit := createTestHabit()

		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, created.ID)
		assert.Equal(t, *habit.Name, *created.Name)
		assert.Equal(t, *habit.Emoji, *created.Emoji)
		assert.Equal(t, habit.UserID, created.UserID)
		assert.Equal(t, *habit.Frequency, *created.Frequency)
		assert.Equal(t, habit.StartDate.Time().Format(time.RFC3339)[:10], created.StartDate.Time().Format(time.RFC3339)[:10])
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})
}

func TestHabitRepository_GetByID(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful get habit", func(t *testing.T) {
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, *habit.Name, *found.Name)
		assert.Equal(t, *habit.Emoji, *found.Emoji)
		assert.Equal(t, *habit.Frequency, *found.Frequency)
	})

	t.Run("habit not found", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), primitive.NewObjectID())
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestHabitRepository_GetAll(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful get all habits for a user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create test habits for the user
		habit1 := createTestHabit()
		habit1.UserID = userID
		name1 := "Habit 1"
		habit1.Name = &name1

		habit2 := createTestHabit()
		habit2.UserID = userID
		name2 := "Habit 2"
		habit2.Name = &name2

		// Create one habit for another user
		otherHabit := createTestHabit()
		otherHabit.UserID = primitive.NewObjectID()

		_, err := repo.Create(context.Background(), habit1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), habit2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), otherHabit)
		require.NoError(t, err)

		// Get habits for the user
		habits, err := repo.GetAll(context.Background(), &userID)
		require.NoError(t, err)
		assert.Len(t, habits, 2)

		// Verify the habit names
		var names []string
		for _, h := range habits {
			names = append(names, *h.Name)
		}
		assert.Contains(t, names, "Habit 1")
		assert.Contains(t, names, "Habit 2")

		// Get all habits (no user filter)
		allHabits, err := repo.GetAll(context.Background(), nil)
		require.NoError(t, err)
		assert.Len(t, allHabits, 3)
	})
}

func TestHabitRepository_Update(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful update habit", func(t *testing.T) {
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		updatedName := "Updated Name"
		newEmoji := "🚀"
		created.Name = &updatedName
		created.Emoji = &newEmoji

		// Wait a moment to ensure timestamp will be different
		time.Sleep(time.Millisecond)

		updated, err := repo.Update(context.Background(), created)
		require.NoError(t, err)
		assert.Equal(t, updatedName, *updated.Name)
		assert.Equal(t, newEmoji, *updated.Emoji)

		// Verify update is persisted
		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedName, *found.Name)
		assert.Equal(t, newEmoji, *found.Emoji)
	})
}

func TestHabitRepository_Delete(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful delete habit", func(t *testing.T) {
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), created.ID)
		require.NoError(t, err)

		// Verify deletion
		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestHabitRepository_AddEntry(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful add habit entry", func(t *testing.T) {
		// First create a habit
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		// Then create an entry for this habit
		entry := createTestHabitEntry()
		entry.HabitID = created.ID
		entry.UserID = created.UserID

		addedEntry, err := repo.AddEntry(context.Background(), entry)
		require.NoError(t, err)
		assert.NotEqual(t, primitive.NilObjectID, addedEntry.ID)
		assert.Equal(t, created.ID, addedEntry.HabitID)
		assert.Equal(t, entry.EntryDate, addedEntry.EntryDate)
		assert.NotEmpty(t, addedEntry.CreatedAt)
		assert.NotEmpty(t, addedEntry.UpdatedAt)
	})
}

func TestHabitRepository_GetEntriesByHabitID(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful get entries for habit", func(t *testing.T) {
		// First create a habit
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		// Add multiple entries
		entry1 := createTestHabitEntry()
		entry1.HabitID = created.ID
		entry1.UserID = created.UserID

		entry2 := createTestHabitEntry()
		entry2.HabitID = created.ID
		entry2.UserID = created.UserID

		addedEntry1, err := repo.AddEntry(context.Background(), entry1)
		require.NoError(t, err)

		addedEntry2, err := repo.AddEntry(context.Background(), entry2)
		require.NoError(t, err)

		// Get all entries for the habit
		entries, err := repo.GetEntriesByHabitID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Len(t, entries, 2)

		// Verify entry IDs
		entryIDs := []primitive.ObjectID{addedEntry1.ID, addedEntry2.ID}
		for _, entry := range entries {
			assert.Contains(t, entryIDs, entry.ID)
		}
	})

	t.Run("no entries for habit", func(t *testing.T) {
		entries, err := repo.GetEntriesByHabitID(context.Background(), primitive.NewObjectID())
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestHabitRepository_UpdateEntry(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful update habit entry", func(t *testing.T) {
		// First create a habit
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		// Add an entry
		entry := createTestHabitEntry()
		entry.HabitID = created.ID
		entry.UserID = created.UserID

		addedEntry, err := repo.AddEntry(context.Background(), entry)
		require.NoError(t, err)

		// Update the entry with new date
		tomorrow := primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour))
		addedEntry.EntryDate = tomorrow

		// Wait a moment to ensure timestamp will be different
		time.Sleep(time.Millisecond)

		updatedEntry, err := repo.UpdateEntry(context.Background(), addedEntry)
		require.NoError(t, err)
		assert.Equal(t, tomorrow, updatedEntry.EntryDate)

		// Verify update is persisted by fetching all entries
		entries, err := repo.GetEntriesByHabitID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
		assert.Equal(t, tomorrow.Time().Format(time.RFC3339)[:10], entries[0].EntryDate.Time().Format(time.RFC3339)[:10])
	})
}

func TestHabitRepository_DeleteEntry(t *testing.T) {
	repo, cleanup := setupHabitTest(t)
	defer cleanup()

	t.Run("successful delete habit entry", func(t *testing.T) {
		// First create a habit
		habit := createTestHabit()
		created, err := repo.Create(context.Background(), habit)
		require.NoError(t, err)

		// Add an entry
		entry := createTestHabitEntry()
		entry.HabitID = created.ID
		entry.UserID = created.UserID

		addedEntry, err := repo.AddEntry(context.Background(), entry)
		require.NoError(t, err)

		// Delete the entry
		err = repo.DeleteEntry(context.Background(), addedEntry.ID)
		require.NoError(t, err)

		// Verify deletion by fetching all entries
		entries, err := repo.GetEntriesByHabitID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}
