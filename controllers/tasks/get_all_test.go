// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/get_all_test.go
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

func TestGetAllTasks(t *testing.T) {
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful get all tasks - with tasks", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Create test tasks
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userID

		task2 := createTestTask()
		task2.ID = primitive.NewObjectID().Hex()
		task2.User = userID

		tasks := []*models.TaskEntity{task1, task2}

		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(tasks, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, task1.ID, response[0].ID)
		assert.Equal(t, task2.ID, response[1].ID)
	})

	t.Run("successful get all tasks - no tasks", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Empty tasks list
		var tasks []*models.TaskEntity

		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(tasks, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(nil, assert.AnError).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
