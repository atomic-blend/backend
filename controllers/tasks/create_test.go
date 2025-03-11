package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTask(t *testing.T) {
	_, mockRepo := setupTest()

	t.Run("successful create task", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID // This should be overwritten by the handler

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(task, nil)

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, response.Title)
		assert.Equal(t, task.StartDate, response.StartDate)
		assert.Equal(t, task.EndDate, response.EndDate)
		assert.Equal(t, userID, response.User) // Verify the task is owned by the authenticated user
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
		controller := NewTaskController(mockRepo)
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
		controller := NewTaskController(mockRepo)
		controller.CreateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
