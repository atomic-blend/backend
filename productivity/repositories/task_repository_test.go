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
