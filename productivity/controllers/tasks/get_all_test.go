// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/get_all_test.go
package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"testing"

	"github.com/atomic-blend/backend/productivity/models"

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
		}), mock.Anything, mock.Anything).Return(tasks, int64(2), nil).Once()

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
		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Tasks, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, task1.ID, response.Tasks[0].ID)
		assert.Equal(t, task2.ID, response.Tasks[1].ID)
	})

	t.Run("successful get all tasks - no tasks", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Empty tasks list
		var tasks []*models.TaskEntity

		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		}), mock.Anything, mock.Anything).Return(tasks, int64(0), nil).Once()

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
		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Tasks, 0)
		assert.Equal(t, int64(0), response.TotalCount)
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
		}), mock.Anything, mock.Anything).Return(nil, int64(0), assert.AnError).Once()

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

	t.Run("successful get all tasks with pagination", func(t *testing.T) {
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
		}), mock.Anything, mock.Anything).Return(tasks, int64(2), nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks?page=1&limit=2", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Tasks, 2)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, int64(1), response.Page)
		assert.Equal(t, int64(2), response.Size)
		assert.Equal(t, int64(1), response.TotalPages)
	})

	t.Run("invalid pagination parameters", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks?page=invalid&limit=2", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("negative pagination parameters", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks?page=-1&limit=2", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
