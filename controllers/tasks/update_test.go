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
	router, mockRepo := setupTest()

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

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)
		mockRepo.On("Update", mock.Anything, taskID, mock.AnythingOfType("*models.TaskEntity")).Return(updatedTask, nil)

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
		controller := NewTaskController(mockRepo)
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

	t.Run("task not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		nonExistentID := primitive.NewObjectID().Hex()
		task := createTestTask()

		mockRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil)

		taskJSON, _ := json.Marshal(task)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+nonExistentID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID}}

		// Call the handler directly
		controller := NewTaskController(mockRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		taskOwnerID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		existingTask := createTestTask()
		existingTask.ID = taskID
		existingTask.User = taskOwnerID // Set a different user as owner

		updatedTask := createTestTask()
		updatedTask.ID = taskID
		updatedTask.Title = "Updated Task"

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)

		taskJSON, _ := json.Marshal(updatedTask)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer(taskJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the handler directly
		controller := NewTaskController(mockRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		existingTask := createTestTask()
		existingTask.ID = taskID
		existingTask.User = userID

		mockRepo.On("GetByID", mock.Anything, taskID).Return(existingTask, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID, bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: taskID}}

		// Call the handler directly
		controller := NewTaskController(mockRepo)
		controller.UpdateTask(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
