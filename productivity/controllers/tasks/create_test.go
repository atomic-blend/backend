// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/create_test.go
package tasks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful create task", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()

		// Set task user and update tag UserIDs to match the authenticated user
		task.User = userID
		if task.Tags != nil && len(*task.Tags) > 0 {
			for _, tag := range *task.Tags {
				tag.UserID = &userID // Set all tags to be owned by the authenticated user
			}
		}

		// Mock tag repository to return valid tags
		if task.Tags != nil && len(*task.Tags) > 0 {
			for _, tag := range *task.Tags {
				// Mock the GetByID call that happens in validation
				mockTagRepo.On("GetByID", mock.Anything, *tag.ID).Return(tag, nil).Once()
			}
		}

		mockTaskRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(task, nil)

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
		assert.Equal(t, userID, response.User) // Verify the task is owned by the authenticated user
		assert.NotNil(t, response.Reminders)
		assert.Len(t, response.Reminders, 2) // Verify reminders were preserved
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		task := createTestTask()
		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Invalid JSON
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid tag - tag not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID

		// Configure mock to return nil for tag lookup, simulating a non-existent tag
		if task.Tags != nil && len(*task.Tags) > 0 {
			tagID := *(*task.Tags)[0].ID // Extract the tag ID from the Tag object
			mockTagRepo.On("GetByID", mock.Anything, tagID).Return(nil, nil).Once()
		}

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Tag not found")
	})

	t.Run("invalid tag - tag belongs to another user", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		anotherUserID := primitive.NewObjectID() // Different user ID
		task := createTestTask()
		task.User = userID

		// Configure mock to return a tag owned by another user
		if task.Tags != nil && len(*task.Tags) > 0 {
			tagID := *(*task.Tags)[0].ID // Extract the tag ID from the Tag object
			tag := &models.Tag{
				ID:     &tagID,
				UserID: &anotherUserID, // Tag belongs to another user
				Name:   "Test Tag",
			}
			mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		}

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "You don't have permission to use this tag")
	})
}
