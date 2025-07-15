// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/delete_test.go
package tasks

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteTask(t *testing.T) {
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful delete task", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		// Ensure task has a properly set User field
		task := createTestTask()
		task.ID = taskID
		task.User = userID // This needs to be a valid ObjectID

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)
		mockTaskRepo.On("Delete", mock.Anything, taskID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing task ID", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Empty task ID
		ctx.Params = []gin.Param{{Key: "id", Value: ""}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("task not found - nil task", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("task not found - error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, errors.New("task not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		taskOwnerID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		task := createTestTask()
		task.ID = taskID
		task.User = taskOwnerID // Set a different user as owner

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("internal server error on delete", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		task := createTestTask()
		task.ID = taskID
		task.User = userID

		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)
		mockTaskRepo.On("Delete", mock.Anything, taskID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.DeleteTask(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
