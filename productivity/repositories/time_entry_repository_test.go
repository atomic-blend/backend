package repositories

import (
	"atomic-blend/backend/productivity/models"
	"atomic-blend/backend/productivity/tests/utils/inmemorymongo"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTimeEntryTest(t *testing.T) (TimeEntryRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewTimeEntryRepository(db)

	// Return cleanup function
	cleanup := func() {
		_ = client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestTimeEntryRepository_Create(t *testing.T) {
	repo, cleanup := setupTimeEntryTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: "2025-05-28T10:00:00Z",
	}

	createdEntry, err := repo.Create(context.Background(), timeEntry)

	assert.NoError(t, err)
	assert.NotNil(t, createdEntry)
	assert.NotNil(t, createdEntry.ID)
	assert.Equal(t, userID, *createdEntry.User)
	assert.Equal(t, timeEntry.StartDate, createdEntry.StartDate)
	assert.Equal(t, timeEntry.EndDate, createdEntry.EndDate)
}

func TestTimeEntryRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTimeEntryTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: "2025-05-28T10:00:00Z",
	}

	createdEntry, err := repo.Create(context.Background(), timeEntry)
	assert.NoError(t, err)

	foundEntry, err := repo.GetByID(context.Background(), createdEntry.ID.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, foundEntry)
	assert.Equal(t, createdEntry.ID, foundEntry.ID)
	assert.Equal(t, userID, *foundEntry.User)
}

func TestTimeEntryRepository_GetAll(t *testing.T) {
	repo, cleanup := setupTimeEntryTest(t)
	defer cleanup()

	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create time entries for user1
	timeEntry1 := &models.TimeEntry{
		User:      &userID1,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: "2025-05-28T10:00:00Z",
	}

	timeEntry2 := &models.TimeEntry{
		User:      &userID1,
		StartDate: "2025-05-28T14:00:00Z",
		EndDate:   "2025-05-28T16:00:00Z",
		Timer:     &[]bool{false}[0],
		Pomodoro:  &[]bool{true}[0],
		CreatedAt: "2025-05-28T14:00:00Z",
		UpdatedAt: "2025-05-28T14:00:00Z",
	}

	// Create time entry for user2
	timeEntry3 := &models.TimeEntry{
		User:      &userID2,
		StartDate: "2025-05-28T09:00:00Z",
		EndDate:   "2025-05-28T11:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T09:00:00Z",
		UpdatedAt: "2025-05-28T09:00:00Z",
	}

	_, err := repo.Create(context.Background(), timeEntry1)
	assert.NoError(t, err)
	_, err = repo.Create(context.Background(), timeEntry2)
	assert.NoError(t, err)
	_, err = repo.Create(context.Background(), timeEntry3)
	assert.NoError(t, err)

	// Get entries for user1
	entries, err := repo.GetAll(context.Background(), &userID1)

	assert.NoError(t, err)
	assert.Len(t, entries, 2)

	// Verify all entries belong to user1
	for _, entry := range entries {
		assert.Equal(t, userID1, *entry.User)
	}
}

func TestTimeEntryRepository_Update(t *testing.T) {
	repo, cleanup := setupTimeEntryTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: "2025-05-28T10:00:00Z",
	}

	createdEntry, err := repo.Create(context.Background(), timeEntry)
	assert.NoError(t, err)

	// Update the entry
	updatedTimeEntry := &models.TimeEntry{
		ID:        createdEntry.ID,
		User:      &userID,
		StartDate: "2025-05-28T11:00:00Z",
		EndDate:   "2025-05-28T13:00:00Z",
		Timer:     &[]bool{false}[0],
		Pomodoro:  &[]bool{true}[0],
		CreatedAt: createdEntry.CreatedAt,
		UpdatedAt: "2025-05-28T11:00:00Z",
	}

	result, err := repo.Update(context.Background(), createdEntry.ID.Hex(), updatedTimeEntry)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedTimeEntry.StartDate, result.StartDate)
	assert.Equal(t, updatedTimeEntry.EndDate, result.EndDate)
	assert.Equal(t, *updatedTimeEntry.Timer, *result.Timer)
	assert.Equal(t, *updatedTimeEntry.Pomodoro, *result.Pomodoro)
}

func TestTimeEntryRepository_Delete(t *testing.T) {
	repo, cleanup := setupTimeEntryTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: "2025-05-28T10:00:00Z",
	}

	createdEntry, err := repo.Create(context.Background(), timeEntry)
	assert.NoError(t, err)

	// Delete the entry
	err = repo.Delete(context.Background(), createdEntry.ID.Hex())
	assert.NoError(t, err)

	// Try to find the deleted entry
	_, err = repo.GetByID(context.Background(), createdEntry.ID.Hex())
	assert.Error(t, err)
}
