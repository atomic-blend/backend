// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/update_test.go
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

func TestUpdateTask(t *testing.T) {
	router, mockTaskRepo, mockTagRepo := setupTest()

	// Set a mock secret for testing
	originalSecret := os.Getenv("SSO_SECRET")
	os.Setenv("SSO_SECRET", "test_secret_for_update_task")
	defer func() {
		os.Setenv("SSO_SECRET", originalSecret)
	}()

	t.Run("successful update task", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		existingTask := createTestTask()
		existingTask.ID = taskID
		existingTask.User = userID

		updatedTask := createTestTask()
		updatedTask.ID = taskID
		updatedTask.User = userID
		updatedTask.Title = "Updated Task"

		// Update reminders for testing
		reminder1 := primitive.NewDateTimeFromTime(time.Now().Add(6 * time.Hour))
		reminder2 := primitive.NewDateTimeFromTime(time.Now().Add(8 * time.Hour))
		reminder3 := primitive.NewDateTimeFromTime(time.Now().Add(10 * time.Hour))
		updatedTask.Reminders = []*primitive.DateTime{&reminder1, &reminder2, &reminder3}

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
		updatedTask.Tags = &tags

		prio := 1
		updatedTask.Priority = &prio

		// Mock tag validation - for each tag, return a valid tag owned by the user
		mockTagRepo.On("GetByID", mock.Anything, tagID1).Return(tags[0], nil).Once()
		mockTagRepo.On("GetByID", mock.Anything, tagID2).Return(tags[1], nil).Once()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
		mockTaskRepo.On("Update", mock.Anything, taskID, mock.AnythingOfType("*models.TaskEntity")).Return(updatedTask, nil)

		taskJSON, _ := json.Marshal(updatedTask)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updatedTask.Title, response.Title)
		assert.Equal(t, updatedTask.StartDate, response.StartDate)
		assert.Equal(t, updatedTask.EndDate, response.EndDate)
		assert.Equal(t, userID, response.User) // Verify the task owner hasn't changed
		assert.NotNil(t, response.Reminders)
		assert.Len(t, response.Reminders, 3) // Verify updated reminders are included
		assert.NotNil(t, response.Tags)
		assert.NotNil(t, response.Priority)
		assert.Len(t, *response.Tags, 2) // Verify tags are included
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		task := createTestTask()
		taskJSON, _ := json.Marshal(task)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new engine with auth middleware but don't set auth user
		engine := gin.New()
		engine.Use(auth.Middleware())
		engine.PUT("/tasks/:id", func(c *gin.Context) {
			router.HandleContext(c)
		})

		// Serve the request without auth header
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("tag validation - invalid tag ID", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()
		invalidTagID := primitive.NewObjectID()

		existingTask := createTestTask()
		existingTask.ID = taskID
		existingTask.User = userID

		updatedTask := createTestTask()
		updatedTask.ID = taskID
		updatedTask.User = userID
		updatedTask.Title = "Updated Task with Invalid Tag"

		// Set up an invalid tag
		tags := []*models.Tag{
			{
				ID:     &invalidTagID,
				UserID: &userID,
				Name:   "Invalid Tag",
			},
		}
		updatedTask.Tags = &tags

		// Mock tag validation to return nil for the invalid tag ID
		mockTagRepo.On("GetByID", mock.Anything, invalidTagID).Return(nil, nil).Once()
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)

		taskJSON, _ := json.Marshal(updatedTask)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Tag not found")
	})

	t.Run("tag validation - tag belongs to another user", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		anotherUserID := primitive.NewObjectID() // Different user ID
		taskID := primitive.NewObjectID().Hex()
		tagID := primitive.NewObjectID()

		existingTask := createTestTask()
		existingTask.ID = taskID
		existingTask.User = userID

		updatedTask := createTestTask()
		updatedTask.ID = taskID
		updatedTask.User = userID
		updatedTask.Title = "Updated Task with Tag from Another User"

		// Set up a tag that belongs to another user
		tags := []*models.Tag{
			{
				ID:     &tagID,
				UserID: &userID, // Initially with correct user ID, but DB will return a different owner
				Name:   "Tag From Another User",
			},
		}
		updatedTask.Tags = &tags

		// Mock tag validation to return a tag owned by another user
		dbTag := &models.Tag{
			ID:     &tagID,
			UserID: &anotherUserID, // Tag belongs to another user
			Name:   "Test Tag",
		}
		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(dbTag, nil).Once()
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)

		taskJSON, _ := json.Marshal(updatedTask)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the handler directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "You don't have permission to use this tag")
	})
}
