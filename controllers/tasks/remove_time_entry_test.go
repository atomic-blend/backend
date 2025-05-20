package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRemoveTimeEntry(t *testing.T) {
	_, mockTaskRepo, _ := setupTest()

	t.Run("successful remove time entry", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID
		task.ID = primitive.NewObjectID().Hex() // Make sure task has a valid ID

		// Create a time entry
		timeEntryID := "test-time-entry-id"
		startDate := time.Now().Format(time.RFC3339)
		endDate := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Add time entry to task
		taskWithTimeEntry := *task
		taskWithTimeEntry.TimeEntries = []*models.TimeEntry{timeEntry}

		// Task after time entry removal (no time entries)
		taskAfterRemoval := *task
		taskAfterRemoval.TimeEntries = []*models.TimeEntry{}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(&taskWithTimeEntry, nil).Once()
		mockTaskRepo.On("RemoveTimeEntry", mock.Anything, task.ID, timeEntryID).Return(&taskAfterRemoval, nil).Once()

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+task.ID+"/time-entries/"+timeEntryID, nil)

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.DELETE("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.RemoveTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()
		timeEntryID := "test-time-entry-id"

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil).Once()

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+taskID+"/time-entries/"+timeEntryID, nil)

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.DELETE("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.RemoveTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("unauthorized - not task owner", func(t *testing.T) {
		// Create authenticated user and different task owner
		userID := primitive.NewObjectID()
		differentUserID := primitive.NewObjectID()
		task := createTestTask()
		task.User = differentUserID             // Task owned by different user
		task.ID = primitive.NewObjectID().Hex() // Make sure task has a valid ID
		timeEntryID := "test-time-entry-id"

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+task.ID+"/time-entries/"+timeEntryID, nil)

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.DELETE("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.RemoveTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("time entry not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID
		task.ID = primitive.NewObjectID().Hex() // Make sure task has a valid ID
		timeEntryID := "non-existent-time-entry-id"

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()
		mockTaskRepo.On("RemoveTimeEntry", mock.Anything, task.ID, timeEntryID).
			Return(nil, errors.New("no time entries found")).Once()

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+task.ID+"/time-entries/"+timeEntryID, nil)

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.DELETE("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.RemoveTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusNotFound, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("database error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID
		task.ID = primitive.NewObjectID().Hex() // Make sure task has a valid ID
		timeEntryID := "test-time-entry-id"

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()
		mockTaskRepo.On("RemoveTimeEntry", mock.Anything, task.ID, timeEntryID).
			Return(nil, errors.New("database error")).Once()

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tasks/"+task.ID+"/time-entries/"+timeEntryID, nil)

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.DELETE("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.RemoveTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})
}
