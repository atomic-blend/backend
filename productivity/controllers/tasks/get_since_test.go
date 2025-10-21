package tasks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTasksSince(t *testing.T) {
	_, mockTaskRepo, _ := setupTest()

	t.Run("successful get tasks since", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tasks := []*models.TaskEntity{
			createTestTask(),
			createTestTask(),
		}
		totalCount := int64(2)

		// Mock repository response
		mockTaskRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(tasks, totalCount, nil).Once()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=2024-01-01T00:00:00Z&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, int64(1), response.Page)
		assert.Equal(t, int64(10), response.Size)
		assert.Equal(t, int64(1), response.TotalPages)
		assert.Len(t, response.Tasks, 2)

		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("missing since parameter", func(t *testing.T) {
		userID := primitive.NewObjectID()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required parameter: since", response["error"])
	})

	t.Run("invalid date format", func(t *testing.T) {
		userID := primitive.NewObjectID()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid date format")
	})

	t.Run("unauthorized request", func(t *testing.T) {
		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", controller.GetTasksSince)

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=2024-01-01T00:00:00Z", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("repository error", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Mock repository error
		mockTaskRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(nil, int64(0), assert.AnError).Once()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=2024-01-01T00:00:00Z&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "assert.AnError")

		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("successful get tasks since with timezone offset", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tasks := []*models.TaskEntity{
			createTestTask(),
		}
		totalCount := int64(1)

		// Mock repository response
		mockTaskRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(tasks, totalCount, nil).Once()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=2024-01-01T12:30:45%2B02:00&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.TotalCount)
		assert.Len(t, response.Tasks, 1)

		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("successful get tasks since without pagination", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tasks := []*models.TaskEntity{
			createTestTask(),
			createTestTask(),
		}
		totalCount := int64(2)

		// Mock repository response - when no pagination, page and limit are nil
		mockTaskRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), (*int64)(nil), (*int64)(nil)).Return(tasks, totalCount, nil).Once()

		controller := NewTaskController(mockTaskRepo, nil)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/tasks/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetTasksSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/tasks/since?since=2024-01-01T00:00:00Z", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedTaskResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, int64(0), response.Page)
		assert.Equal(t, int64(0), response.Size)
		assert.Equal(t, int64(0), response.TotalPages)
		assert.Len(t, response.Tasks, 2)

		mockTaskRepo.AssertExpectations(t)
	})
}
