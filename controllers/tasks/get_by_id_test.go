package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTaskByID(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get task by id", func(t *testing.T) {
		// Create authenticated user
		userId := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()
		task := createTestTask()
		task.ID = taskID
		task.User = userId // Set the task owner

		mockRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userId})
		// Copy headers and params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
	})

	t.Run("task not found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID().Hex()
		userId := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+nonExistentID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userId})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID}}

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Task not found", response["error"])
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		taskID := primitive.NewObjectID().Hex()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserId := primitive.NewObjectID()
		taskOwnerId := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		task := createTestTask()
		task.ID = taskID
		task.User = taskOwnerId // Set a different user as owner

		mockRepo.On("GetByID", mock.Anything, taskID).Return(task, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks/"+taskID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserId})
		// Copy headers and params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.GetTaskByID(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
