package tasks

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
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

func TestAddTimeEntry(t *testing.T) {
	_, mockTaskRepo, mockTagRepo := setupTest()

	t.Run("successful add time entry", func(t *testing.T) {
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

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()

		// Create a task with the time entry added
		taskWithTimeEntry := *task
		taskWithTimeEntry.TimeEntries = []*models.TimeEntry{timeEntry}
		mockTaskRepo.On("AddTimeEntry", mock.Anything, task.ID, mock.AnythingOfType("*models.TimeEntry")).Return(&taskWithTimeEntry, nil).Once()

		// Create request body
		timeEntryJSON, _ := json.Marshal(timeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/"+task.ID+"/time-entries", bytes.NewBuffer(timeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		router.POST("/tasks/:id/time-entries", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.AddTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusOK, w.Code)

		// Verify response body
		var response models.TaskEntity
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, task.ID, response.ID)
		assert.Len(t, response.TimeEntries, 1)
		assert.Equal(t, timeEntryID, *response.TimeEntries[0].ID)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()

		// Create a time entry
		timeEntryID := "test-time-entry-id"
		startDate := time.Now().Format(time.RFC3339)
		endDate := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil).Once()

		// Create request body
		timeEntryJSON, _ := json.Marshal(timeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/"+taskID+"/time-entries", bytes.NewBuffer(timeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		router.POST("/tasks/:id/time-entries", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.AddTimeEntry(c)
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

		// Create a time entry
		timeEntryID := "test-time-entry-id"
		startDate := time.Now().Format(time.RFC3339)
		endDate := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
		timeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()

		// Create request body
		timeEntryJSON, _ := json.Marshal(timeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/"+task.ID+"/time-entries", bytes.NewBuffer(timeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		router.POST("/tasks/:id/time-entries", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.AddTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("error adding time entry", func(t *testing.T) {
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

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()
		mockTaskRepo.On("AddTimeEntry", mock.Anything, task.ID, mock.AnythingOfType("*models.TimeEntry")).
			Return(nil, errors.New("database error")).Once()

		// Create request body
		timeEntryJSON, _ := json.Marshal(timeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tasks/"+task.ID+"/time-entries", bytes.NewBuffer(timeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, mockTagRepo)
		router.POST("/tasks/:id/time-entries", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.AddTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})
}
