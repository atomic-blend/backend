package repositories

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/utils/inmemorymongo"
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

func TestTaskRepository_BulkUpdate(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()

	t.Run("bulk update with no conflicts", func(t *testing.T) {
		// Create two tasks first
		task1 := &models.TaskEntity{
			Title: "Task 1",
			User:  userID,
		}
		task2 := &models.TaskEntity{
			Title: "Task 2",
			User:  userID,
		}

		createdTask1, err := repo.Create(context.Background(), task1)
		require.NoError(t, err)
		createdTask2, err := repo.Create(context.Background(), task2)
		require.NoError(t, err)

		// Update both tasks
		updateTime := time.Now().Add(1 * time.Hour)
		updatedTask1 := &models.TaskEntity{
			ID:        createdTask1.ID,
			Title:     "Updated Task 1",
			User:      userID,
			UpdatedAt: primitive.NewDateTimeFromTime(updateTime),
		}
		updatedTask2 := &models.TaskEntity{
			ID:        createdTask2.ID,
			Title:     "Updated Task 2",
			User:      userID,
			UpdatedAt: primitive.NewDateTimeFromTime(updateTime),
		}

		// Perform bulk update
		updated, conflicts, err := repo.BulkUpdate(context.Background(), []*models.TaskEntity{updatedTask1, updatedTask2})
		require.NoError(t, err)

		// Should update both tasks with no conflicts
		assert.Len(t, updated, 2)
		assert.Len(t, conflicts, 0)
		assert.Equal(t, "Updated Task 1", updated[0].Title)
		assert.Equal(t, "Updated Task 2", updated[1].Title)
	})

	t.Run("bulk update with conflicts", func(t *testing.T) {
		// Create a task
		task := &models.TaskEntity{
			Title: "Original Task",
			User:  userID,
		}

		createdTask, err := repo.Create(context.Background(), task)
		require.NoError(t, err)

		// Update the task in the database first (making it more recent)
		recentUpdate := &models.TaskEntity{
			ID:    createdTask.ID,
			Title: "Recently Updated Task",
			User:  userID,
		}
		_, err = repo.Update(context.Background(), createdTask.ID, recentUpdate)
		require.NoError(t, err)

		// Try to bulk update with an older timestamp (should create conflict)
		olderUpdateTime := time.Now().Add(-1 * time.Hour)
		conflictingTask := &models.TaskEntity{
			ID:        createdTask.ID,
			Title:     "Conflicting Update",
			User:      userID,
			UpdatedAt: primitive.NewDateTimeFromTime(olderUpdateTime),
		}

		// Perform bulk update
		updated, conflicts, err := repo.BulkUpdate(context.Background(), []*models.TaskEntity{conflictingTask})
		require.NoError(t, err)

		// Should have no updates and one conflict
		assert.Len(t, updated, 0)
		assert.Len(t, conflicts, 1)
		assert.Equal(t, "task", conflicts[0].Type)
		assert.Equal(t, "Recently Updated Task", conflicts[0].OldItem.(*models.TaskEntity).Title)
		assert.Equal(t, "Conflicting Update", conflicts[0].NewItem.(*models.TaskEntity).Title)
	})

	t.Run("bulk update with new task creation", func(t *testing.T) {
		// Try to update a task that doesn't exist (should create it)
		newTaskID := primitive.NewObjectID().Hex()
		newTask := &models.TaskEntity{
			ID:    newTaskID,
			Title: "New Task",
			User:  userID,
		}

		// Perform bulk update
		updated, conflicts, err := repo.BulkUpdate(context.Background(), []*models.TaskEntity{newTask})
		require.NoError(t, err)

		// Should create the new task
		assert.Len(t, updated, 1)
		assert.Len(t, conflicts, 0)
		assert.Equal(t, "New Task", updated[0].Title)
		assert.Equal(t, newTaskID, updated[0].ID)
	})

	t.Run("bulk update mixed scenario", func(t *testing.T) {
		// Create an existing task
		existingTask := &models.TaskEntity{
			Title: "Existing Task",
			User:  userID,
		}
		createdTask, err := repo.Create(context.Background(), existingTask)
		require.NoError(t, err)

		// Update it in the database first
		_, err = repo.Update(context.Background(), createdTask.ID, &models.TaskEntity{
			ID:    createdTask.ID,
			Title: "Database Updated Task",
			User:  userID,
		})
		require.NoError(t, err)

		// Prepare bulk update with:
		// 1. A conflicting update (older timestamp)
		// 2. A new task that doesn't exist
		// 3. A successful update (newer timestamp)

		olderUpdateTime := time.Now().Add(-1 * time.Hour)
		newerUpdateTime := time.Now().Add(1 * time.Hour)

		conflictingTask := &models.TaskEntity{
			ID:        createdTask.ID,
			Title:     "Conflicting Update",
			User:      userID,
			UpdatedAt: primitive.NewDateTimeFromTime(olderUpdateTime),
		}

		newTaskID := primitive.NewObjectID().Hex()
		newTask := &models.TaskEntity{
			ID:    newTaskID,
			Title: "Brand New Task",
			User:  userID,
		}

		// Create another existing task for successful update
		anotherTask := &models.TaskEntity{
			Title: "Another Task",
			User:  userID,
		}
		createdAnotherTask, err := repo.Create(context.Background(), anotherTask)
		require.NoError(t, err)

		successfulUpdate := &models.TaskEntity{
			ID:        createdAnotherTask.ID,
			Title:     "Successfully Updated Task",
			User:      userID,
			UpdatedAt: primitive.NewDateTimeFromTime(newerUpdateTime),
		}

		// Perform bulk update
		updated, conflicts, err := repo.BulkUpdate(context.Background(), []*models.TaskEntity{
			conflictingTask,
			newTask,
			successfulUpdate,
		})
		require.NoError(t, err)

		// Should have 2 updates (new task + successful update) and 1 conflict
		assert.Len(t, updated, 2)
		assert.Len(t, conflicts, 1)

		// Check conflict
		assert.Equal(t, "task", conflicts[0].Type)
		assert.Equal(t, "Database Updated Task", conflicts[0].OldItem.(*models.TaskEntity).Title)
		assert.Equal(t, "Conflicting Update", conflicts[0].NewItem.(*models.TaskEntity).Title)

		// Check updates
		updatedTitles := []string{updated[0].Title, updated[1].Title}
		assert.Contains(t, updatedTitles, "Brand New Task")
		assert.Contains(t, updatedTitles, "Successfully Updated Task")
	})
}