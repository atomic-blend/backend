// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/get_by_id_test.go
package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	"atomic-blend/backend/productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTaskByID(t *testing.T) {
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful get task by id", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()
		task := createTestTask()
		task.ID = taskID
		task.User = userID // Set the task owner

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		// Copy headers and params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
		assert.NotNil(t, response.Reminders)
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing task ID", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: ""}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("task not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, assert.AnError)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
