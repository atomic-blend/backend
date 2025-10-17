package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

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

	t.Run("get all tasks without pagination", func(t *testing.T) {
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

		// Test getting all tasks for specific user (no pagination)
		found, totalCount, err := repo.GetAll(context.Background(), &userID, nil, nil)
		require.NoError(t, err)
		assert.Len(t, found, 2)
		assert.Equal(t, int64(2), totalCount)
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

		// Test getting all tasks for all users (no pagination)
		allTasks, totalCount, err := repo.GetAll(context.Background(), nil, nil, nil)
		require.NoError(t, err)
		assert.Len(t, allTasks, 2)
		assert.Equal(t, int64(2), totalCount)
	})

	t.Run("get all tasks with pagination", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create 7 tasks for testing pagination
		for i := 1; i <= 7; i++ {
			task := createTestTask()
			task.Title = fmt.Sprintf("Task %d", i)
			task.User = userID
			_, err := repo.Create(context.Background(), task)
			require.NoError(t, err)
		}

		// Test getting all tasks (no pagination)
		allTasks, totalCount, err := repo.GetAll(context.Background(), &userID, nil, nil)
		require.NoError(t, err)
		assert.Len(t, allTasks, 7)
		assert.Equal(t, int64(7), totalCount)

		// Test pagination - first page with limit 3
		page1 := int64(1)
		limit := int64(3)
		page1Tasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page1, &limit)
		require.NoError(t, err)
		assert.Len(t, page1Tasks, 3)
		assert.Equal(t, int64(7), totalCount)

		// Test pagination - second page with limit 3
		page2 := int64(2)
		page2Tasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page2, &limit)
		require.NoError(t, err)
		assert.Len(t, page2Tasks, 3)
		assert.Equal(t, int64(7), totalCount)

		// Test pagination - third page with limit 3 (should have 1 task)
		page3 := int64(3)
		page3Tasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page3, &limit)
		require.NoError(t, err)
		assert.Len(t, page3Tasks, 1)
		assert.Equal(t, int64(7), totalCount)

		// Test pagination - fourth page with limit 3 (should be empty)
		page4 := int64(4)
		page4Tasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page4, &limit)
		require.NoError(t, err)
		assert.Len(t, page4Tasks, 0)
		assert.Equal(t, int64(7), totalCount)
	})

	t.Run("pagination edge cases", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create 3 tasks for testing edge cases
		for i := 1; i <= 3; i++ {
			task := createTestTask()
			task.Title = fmt.Sprintf("Task %d", i)
			task.User = userID
			_, err := repo.Create(context.Background(), task)
			require.NoError(t, err)
		}

		// Test with only page provided (should return all tasks)
		pageOnly := int64(1)
		pageOnlyTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &pageOnly, nil)
		require.NoError(t, err)
		assert.Len(t, pageOnlyTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Test with only limit provided (should return all tasks)
		limitOnly := int64(2)
		limitOnlyTasks, totalCount, err := repo.GetAll(context.Background(), &userID, nil, &limitOnly)
		require.NoError(t, err)
		assert.Len(t, limitOnlyTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Test with page = 0 (should return all tasks)
		pageZero := int64(0)
		pageZeroTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &pageZero, &limitOnly)
		require.NoError(t, err)
		assert.Len(t, pageZeroTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Test with limit = 0 (should return all tasks)
		limitZero := int64(0)
		limitZeroTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &pageOnly, &limitZero)
		require.NoError(t, err)
		assert.Len(t, limitZeroTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Test with negative page (should return all tasks)
		pageNegative := int64(-1)
		pageNegativeTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &pageNegative, &limitOnly)
		require.NoError(t, err)
		assert.Len(t, pageNegativeTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Test with negative limit (should return all tasks)
		limitNegative := int64(-1)
		limitNegativeTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &pageOnly, &limitNegative)
		require.NoError(t, err)
		assert.Len(t, limitNegativeTasks, 3)
		assert.Equal(t, int64(3), totalCount)
	})

	t.Run("empty result sets", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Test getting tasks for user with no tasks (no pagination)
		emptyTasks, totalCount, err := repo.GetAll(context.Background(), &userID, nil, nil)
		require.NoError(t, err)
		assert.Len(t, emptyTasks, 0)
		assert.Equal(t, int64(0), totalCount)

		// Test pagination with empty result set
		page1 := int64(1)
		limit := int64(10)
		emptyPageTasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page1, &limit)
		require.NoError(t, err)
		assert.Len(t, emptyPageTasks, 0)
		assert.Equal(t, int64(0), totalCount)
	})

	t.Run("sorting by created_at desc", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create tasks with slight delays to ensure different created_at times
		task1 := createTestTask()
		task1.Title = "First Task"
		task1.User = userID
		_, err := repo.Create(context.Background(), task1)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond) // Small delay to ensure different timestamps

		task2 := createTestTask()
		task2.Title = "Second Task"
		task2.User = userID
		_, err = repo.Create(context.Background(), task2)
		require.NoError(t, err)

		time.Sleep(10 * time.Millisecond)

		task3 := createTestTask()
		task3.Title = "Third Task"
		task3.User = userID
		_, err = repo.Create(context.Background(), task3)
		require.NoError(t, err)

		// Test that tasks are returned in descending order by created_at
		allTasks, totalCount, err := repo.GetAll(context.Background(), &userID, nil, nil)
		require.NoError(t, err)
		assert.Len(t, allTasks, 3)
		assert.Equal(t, int64(3), totalCount)

		// Should be ordered by created_at desc (most recent first)
		assert.Equal(t, "Third Task", allTasks[0].Title)
		assert.Equal(t, "Second Task", allTasks[1].Title)
		assert.Equal(t, "First Task", allTasks[2].Title)

		// Test pagination maintains the same order
		page1 := int64(1)
		limit := int64(2)
		page1Tasks, totalCount, err := repo.GetAll(context.Background(), &userID, &page1, &limit)
		require.NoError(t, err)
		assert.Len(t, page1Tasks, 2)
		assert.Equal(t, int64(3), totalCount)
		assert.Equal(t, "Third Task", page1Tasks[0].Title)
		assert.Equal(t, "Second Task", page1Tasks[1].Title)
	})
}

func TestTaskRepository_DeleteByUserID(t *testing.T) {
	repo, cleanup := setupTaskTest(t)
	defer cleanup()

	// Create test users
	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create tasks for user 1
	task1 := createTestTask()
	task1.User = userID1
	createdTask1, err := repo.Create(context.Background(), task1)
	require.NoError(t, err)
	require.NotNil(t, createdTask1)

	task2 := createTestTask()
	task2.Title = "Task 2 for User 1"
	task2.User = userID1
	createdTask2, err := repo.Create(context.Background(), task2)
	require.NoError(t, err)
	require.NotNil(t, createdTask2)

	// Create tasks for user 2
	task3 := createTestTask()
	task3.Title = "Task 1 for User 2"
	task3.User = userID2
	createdTask3, err := repo.Create(context.Background(), task3)
	require.NoError(t, err)
	require.NotNil(t, createdTask3)

	// Verify all tasks exist
	allTasks, totalCount, err := repo.GetAll(context.Background(), nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, allTasks, 3)
	assert.Equal(t, int64(3), totalCount)

	// Delete all tasks for user 1
	err = repo.DeleteByUserID(context.Background(), userID1)
	require.NoError(t, err)

	// Verify only user 2's task remains
	remainingTasks, totalCount, err := repo.GetAll(context.Background(), nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, remainingTasks, 1)
	assert.Equal(t, int64(1), totalCount)
	assert.Equal(t, "Task 1 for User 2", remainingTasks[0].Title)
	assert.Equal(t, userID2, remainingTasks[0].User)

	// Verify user 1's tasks are gone
	user1Tasks, totalCount, err := repo.GetAll(context.Background(), &userID1, nil, nil)
	require.NoError(t, err)
	assert.Len(t, user1Tasks, 0)
	assert.Equal(t, int64(0), totalCount)

	// Delete all tasks for user 2
	err = repo.DeleteByUserID(context.Background(), userID2)
	require.NoError(t, err)

	// Verify no tasks remain
	finalTasks, totalCount, err := repo.GetAll(context.Background(), nil, nil, nil)
	require.NoError(t, err)
	assert.Len(t, finalTasks, 0)
	assert.Equal(t, int64(0), totalCount)
}
