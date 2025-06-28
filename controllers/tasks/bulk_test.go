package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestBulkUpdateTasks(t *testing.T) {
	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_bulk_update_tasks")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("successful bulk update with no conflicts", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		
		// Create test tasks without tags to avoid validation calls
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID
		task1.Title = "Updated Task 1"
		task1.Tags = nil // Remove tags to avoid validation calls
		
		task2 := createTestTask()
		task2.ID = primitive.NewObjectID().Hex()
		task2.User = userID
		task2.Title = "Updated Task 2"
		task2.Tags = nil // Remove tags to avoid validation calls

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1, task2},
		}

		// Mock repository response - no conflicts, all updated
		updatedTasks := []*models.TaskEntity{task1, task2}
		var conflicts []*models.ConflictedItem
		
		mockTaskRepo.On("BulkUpdate", mock.Anything, mock.AnythingOfType("[]*models.TaskEntity")).Return(updatedTasks, conflicts, nil)

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.BulkTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Updated, 2)
		assert.Len(t, response.Conflicts, 0)
		assert.Equal(t, task1.Title, response.Updated[0].Title)
		assert.Equal(t, task2.Title, response.Updated[1].Title)
	})

	t.Run("bulk update with conflicts", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		
		// Create test tasks
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID
		task1.Title = "Updated Task 1"
		task1.Tags = nil // Remove tags to avoid validation calls
		task1.UpdatedAt = primitive.NewDateTimeFromTime(time.Now().Add(-1 * time.Hour)) // Older timestamp

		task2 := createTestTask()
		task2.ID = primitive.NewObjectID().Hex()
		task2.User = userID
		task2.Title = "Updated Task 2"
		task2.Tags = nil // Remove tags to avoid validation calls

		// Create an existing task that's more recent (will cause conflict)
		existingTask1 := createTestTask()
		existingTask1.ID = task1.ID
		existingTask1.User = userID
		existingTask1.Title = "Existing Task 1 (More Recent)"
		existingTask1.Tags = nil // Remove tags to avoid validation calls
		existingTask1.UpdatedAt = primitive.NewDateTimeFromTime(time.Now()) // More recent timestamp

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1, task2},
		}

		// Mock repository response - one conflict, one update
		updatedTasks := []*models.TaskEntity{task2} // Only task2 gets updated
		conflicts := []*models.ConflictedItem{
			{
				Type: "task",
				Old:  existingTask1,
				New:  task1,
			},
		}
		
		mockTaskRepo.On("BulkUpdate", mock.Anything, mock.AnythingOfType("[]*models.TaskEntity")).Return(updatedTasks, conflicts, nil)

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.BulkTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Updated, 1)
		assert.Len(t, response.Conflicts, 1)
		assert.Equal(t, task2.Title, response.Updated[0].Title)
		assert.Equal(t, "task", response.Conflicts[0].Type)
	})

	t.Run("bulk update with tags validation", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		
		// Create test task with tags
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID
		task1.Title = "Task with Tags"

		// Configure task with tags
		tagID1 := primitive.NewObjectID()
		tagID2 := primitive.NewObjectID()
		tags := []*models.Tag{
			{
				ID:     &tagID1,
				UserID: &userID,
				Name:   "Test Tag 1",
			},
			{
				ID:     &tagID2,
				UserID: &userID,
				Name:   "Test Tag 2",
			},
		}
		task1.Tags = &tags

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1},
		}

		// Mock tag validation - for each tag, return a valid tag owned by the user
		mockTagRepo.On("GetByID", mock.Anything, tagID1).Return(tags[0], nil).Once()
		mockTagRepo.On("GetByID", mock.Anything, tagID2).Return(tags[1], nil).Once()

		// Mock repository response
		updatedTasks := []*models.TaskEntity{task1}
		var conflicts []*models.ConflictedItem
		
		mockTaskRepo.On("BulkUpdate", mock.Anything, mock.AnythingOfType("[]*models.TaskEntity")).Return(updatedTasks, conflicts, nil)

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.BulkTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Updated, 1)
		assert.Len(t, response.Conflicts, 0)
		assert.NotNil(t, response.Updated[0].Tags)
		assert.Len(t, *response.Updated[0].Tags, 2)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1},
		}

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new engine with auth middleware but don't set auth user
		engine := gin.New()
		engine.Use(auth.Middleware())
		engine.PUT("/tasks/bulk", func(c *gin.Context) {
			// This will fail because no router is set up
		})

		// Serve the request without auth header
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("empty tasks array", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		
		// Create request with empty tasks array
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{},
		}

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "At least one task is required")
	})

	t.Run("task without ID", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		
		// Create test task without ID
		task1 := createTestTask()
		task1.ID = "" // No ID provided
		task1.User = userID
		task1.Tags = nil // Remove tags to avoid validation calls

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1},
		}

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Task ID is required for bulk update")
		assert.Equal(t, float64(0), response["index"]) // Index of the problematic task
	})

	t.Run("tag validation - invalid tag ID", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		invalidTagID := primitive.NewObjectID()

		// Create test task with invalid tag
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID
		task1.Title = "Task with Invalid Tag"

		// Set up an invalid tag
		tags := []*models.Tag{
			{
				ID:     &invalidTagID,
				UserID: &userID,
				Name:   "Invalid Tag",
			},
		}
		task1.Tags = &tags

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1},
		}

		// Mock tag validation to return nil for the invalid tag ID
		mockTagRepo.On("GetByID", mock.Anything, invalidTagID).Return(nil, nil).Once()

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Tag not found")
		assert.Equal(t, float64(0), response["index"]) // Index of the problematic task
	})

	t.Run("tag validation - tag belongs to another user", func(t *testing.T) {
		// Create new mocks for this test
		_, mockTaskRepo, mockTagRepo := setupTest()
		
		// Create authenticated user
		userID := primitive.NewObjectID()
		anotherUserID := primitive.NewObjectID() // Different user ID
		tagID := primitive.NewObjectID()

		// Create test task with tag from another user
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID
		task1.Title = "Task with Tag from Another User"

		// Set up a tag that belongs to another user
		tags := []*models.Tag{
			{
				ID:     &tagID,
				UserID: &userID, // Initially with correct user ID, but DB will return a different owner
				Name:   "Tag From Another User",
			},
		}
		task1.Tags = &tags

		// Create request
		request := models.BulkTaskRequest{
			Tasks: []*models.TaskEntity{task1},
		}

		// Mock tag validation to return a tag owned by another user
		dbTag := &models.Tag{
			ID:     &tagID,
			UserID: &anotherUserID, // Tag belongs to another user
			Name:   "Test Tag",
		}
		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(dbTag, nil).Once()

		requestJSON, _ := json.Marshal(request)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/bulk", bytes.NewBuffer(requestJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.BulkUpdateTasks(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "You don't have permission to use this tag")
		assert.Equal(t, float64(0), response["index"]) // Index of the problematic task
	})
}
