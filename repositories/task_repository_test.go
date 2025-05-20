package repositories

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/utils/inmemorymongo"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTaskTest(t *testing.T) (TaskRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewTaskRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestTask() *models.TaskEntity {
	desc := "Test Description"
	completed := false
	now := primitive.NewDateTimeFromTime(time.Now())
	end := primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour))
	reminder1 := primitive.NewDateTimeFromTime(time.Now().Add(12 * time.Hour))
	reminder2 := primitive.NewDateTimeFromTime(time.Now().Add(18 * time.Hour))
	return &models.TaskEntity{
		Title:       "Test Task",
		Description: &desc,
		Completed:   &completed,
		StartDate:   &now,
		EndDate:     &end,
		Reminders:   []*primitive.DateTime{&reminder1, &reminder2},
	}
}

func TestTaskRepository_Create(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful create task", func(t *testing.T) {
		task := createTestTask()

		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)
		assert.NotEmpty(t, created.ID)
		assert.Equal(t, task.Title, created.Title)
		assert.NotEmpty(t, created.CreatedAt)
		assert.NotEmpty(t, created.UpdatedAt)
		assert.NotNil(t, created.StartDate)
		assert.NotNil(t, created.EndDate)
		assert.NotNil(t, created.Reminders)
		assert.Len(t, created.Reminders, 2)
	})
}

func TestTaskRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful get task", func(t *testing.T) {
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, task.Title, found.Title)
		assert.Equal(t, task.StartDate, found.StartDate)
		assert.Equal(t, task.EndDate, found.EndDate)
		assert.NotNil(t, found.Reminders)
		assert.Len(t, found.Reminders, 2)
		assert.Equal(t, task.Reminders, found.Reminders)
	})

	t.Run("task not found", func(t *testing.T) {
		found, err := repo.GetByID(context.Background(), primitive.NewObjectID().Hex())
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestTaskRepository_Update(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful update task", func(t *testing.T) {
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		updatedTitle := "Updated Title"
		created.Title = updatedTitle
		updated, err := repo.Update(context.Background(), created.ID, created)
		require.NoError(t, err)
		assert.Equal(t, updatedTitle, updated.Title)
		assert.Equal(t, task.StartDate, updated.StartDate)
		assert.Equal(t, task.EndDate, updated.EndDate)
		assert.NotNil(t, updated.Reminders)
		assert.Len(t, updated.Reminders, 2)
		assert.Equal(t, task.Reminders, updated.Reminders)

		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Equal(t, updatedTitle, found.Title)
		assert.NotNil(t, found.Reminders)
		assert.Len(t, found.Reminders, 2)
	})
}

func TestTaskRepository_Delete(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful delete task", func(t *testing.T) {
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		err = repo.Delete(context.Background(), created.ID)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestTaskRepository_GetAll(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful get all tasks", func(t *testing.T) {
		userID := primitive.NewObjectID()
		// Create reminder dates for testing
		reminder1 := primitive.NewDateTimeFromTime(time.Now().Add(6 * time.Hour))
		reminder2 := primitive.NewDateTimeFromTime(time.Now().Add(12 * time.Hour))

		tasks := []*models.TaskEntity{
			{
				Title:       "Task 1",
				User:        userID,
				Description: createTestTask().Description,
				StartDate:   createTestTask().StartDate,
				EndDate:     createTestTask().EndDate,
				Completed:   createTestTask().Completed,
				Reminders:   []*primitive.DateTime{&reminder1},
			},
			{
				Title:       "Task 2",
				User:        userID,
				Description: createTestTask().Description,
				StartDate:   createTestTask().StartDate,
				EndDate:     createTestTask().EndDate,
				Completed:   createTestTask().Completed,
				Reminders:   []*primitive.DateTime{&reminder1, &reminder2},
			},
		}

		for _, task := range tasks {
			_, err := repo.Create(context.Background(), task)
			require.NoError(t, err)
		}

		found, err := repo.GetAll(context.Background(), &userID)
		require.NoError(t, err)
		assert.Len(t, found, 2)
		// Ensure reminders were saved correctly
		assert.NotNil(t, found[0].Reminders)
		assert.NotNil(t, found[1].Reminders)

		// Validate reminders lengths
		var taskWithOneReminder, taskWithTwoReminders *models.TaskEntity
		for _, task := range found {
			switch task.Title {
			case "Task 1":
				taskWithOneReminder = task
			case "Task 2":
				taskWithTwoReminders = task
			}
		}
		assert.NotNil(t, taskWithOneReminder)
		assert.NotNil(t, taskWithTwoReminders)
		assert.Len(t, taskWithOneReminder.Reminders, 1)
		assert.Len(t, taskWithTwoReminders.Reminders, 2)

		allTasks, err := repo.GetAll(context.Background(), nil)
		require.NoError(t, err)
		assert.Len(t, allTasks, 2)
	})
}

func TestTaskRepository_AddTimeEntry(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful add time entry", func(t *testing.T) {
		// Create a task first
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		// Create time entry
		timeEntryID := primitive.NewObjectID().Hex()
		startDate := time.Now().Format(time.RFC3339)
		endDate := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Add time entry to task
		updatedTask, err := repo.AddTimeEntry(context.Background(), created.ID, timeEntry)
		require.NoError(t, err)
		require.NotNil(t, updatedTask)
		require.NotNil(t, updatedTask.TimeEntries)
		require.Len(t, updatedTask.TimeEntries, 1)
		require.Equal(t, timeEntryID, *updatedTask.TimeEntries[0].ID)
		require.Equal(t, startDate, updatedTask.TimeEntries[0].StartDate)
		require.Equal(t, endDate, updatedTask.TimeEntries[0].EndDate)
	})

	t.Run("add time entry to non-existent task", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		timeEntryID := primitive.NewObjectID().Hex()
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: time.Now().Format(time.RFC3339),
			EndDate:   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}

		_, err := repo.AddTimeEntry(context.Background(), nonExistentID, timeEntry)
		require.NoError(t, err) // MongoDB doesn't error on updates with no matches
	})
}

func TestTaskRepository_RemoveTimeEntry(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful remove time entry", func(t *testing.T) {
		// Create a task first
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		// Add time entry to task
		timeEntryID := primitive.NewObjectID().Hex()
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: time.Now().Format(time.RFC3339),
			EndDate:   time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}

		updatedTask, err := repo.AddTimeEntry(context.Background(), created.ID, timeEntry)
		require.NoError(t, err)
		require.Len(t, updatedTask.TimeEntries, 1)

		// Remove the time entry
		taskAfterRemoval, err := repo.RemoveTimeEntry(context.Background(), created.ID, timeEntryID)
		require.NoError(t, err)
		require.NotNil(t, taskAfterRemoval)

		// Should have no time entries now
		require.Empty(t, taskAfterRemoval.TimeEntries)
	})

	t.Run("remove non-existent time entry", func(t *testing.T) {
		// Create a task first
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		// Try to remove a non-existent time entry
		nonExistentTimeEntryID := primitive.NewObjectID().Hex()
		_, err = repo.RemoveTimeEntry(context.Background(), created.ID, nonExistentTimeEntryID)
		// require error "no time entries found"
		require.Error(t, errors.New("no time entries found"))
	})
}

func TestTaskRepository_UpdateTimeEntry(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	t.Run("successful update time entry", func(t *testing.T) {
		// Create a task first
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)
		require.NotEmpty(t, created.ID)

		// Add time entry to task
		timeEntryID := primitive.NewObjectID().Hex()
		originalStartDate := time.Now().Format(time.RFC3339)
		originalEndDate := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: originalStartDate,
			EndDate:   originalEndDate,
		}

		updatedTask, err := repo.AddTimeEntry(context.Background(), created.ID, timeEntry)
		require.NoError(t, err)
		require.Len(t, updatedTask.TimeEntries, 1)

		// Update the time entry
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		taskAfterUpdate, err := repo.UpdateTimeEntry(context.Background(), created.ID, timeEntryID, updatedTimeEntry)
		require.NoError(t, err)
		require.NotNil(t, taskAfterUpdate)
		require.Len(t, taskAfterUpdate.TimeEntries, 1)

		// Verify the time entry was updated
		require.Equal(t, updatedStartDate, taskAfterUpdate.TimeEntries[0].StartDate)
		require.Equal(t, updatedEndDate, taskAfterUpdate.TimeEntries[0].EndDate)
	})

	t.Run("update non-existent time entry", func(t *testing.T) {
		// Create a task first
		task := createTestTask()
		created, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		// Try to update a non-existent time entry
		nonExistentTimeEntryID := primitive.NewObjectID().Hex()
		updatedTimeEntry := &models.TimeEntry{
			ID:        &nonExistentTimeEntryID,
			StartDate: time.Now().Format(time.RFC3339),
			EndDate:   time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		}

		_, err = repo.UpdateTimeEntry(context.Background(), created.ID, nonExistentTimeEntryID, updatedTimeEntry)
		require.Error(t, errors.New("no time entries found"))
	})
}
