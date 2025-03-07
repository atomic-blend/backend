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
	router, mockRepo := setupTest()

	t.Run("successful get all tasks - with tasks", func(t *testing.T) {
		// Create authenticated user
		userId := primitive.NewObjectID()

		// Create test tasks
		task1 := createTestTask()
		task1.ID = primitive.NewObjectID().Hex()
		task1.User = userId

		task2 := createTestTask()
		task2.ID = primitive.NewObjectID().Hex()
		task2.User = userId

		tasks := []*models.TaskEntity{task1, task2}

		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return userID.Hex() == userId.Hex()
		})).Return(tasks, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userId})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.TaskEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
		assert.Equal(t, task1.Title, response[0].Title)
		assert.Equal(t, task2.Title, response[1].Title)
	})

	t.Run("successful get all tasks - empty list", func(t *testing.T) {
		// Create authenticated user
		userId := primitive.NewObjectID()

		// Return nil to simulate no tasks
		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return userID.Hex() == userId.Hex()
		})).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userId})

		// Call the controller directly with our context that has auth
		controller := NewTaskController(mockRepo)
		controller.GetAllTasks(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check that the response is an empty array, not null
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tasks", nil)

		// Call the endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})
}
