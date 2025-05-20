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

func TestUpdateTimeEntry(t *testing.T) {
	_, mockTaskRepo, _ := setupTest()

	t.Run("successful update time entry", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		task := createTestTask()
		task.User = userID
		task.ID = primitive.NewObjectID().Hex() // Make sure task has a valid ID

		// Create initial time entry
		timeEntryID := primitive.NewObjectID()
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

		// Updated time entry values
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		// Task after time entry update
		taskAfterUpdate := *task
		taskAfterUpdate.TimeEntries = []*models.TimeEntry{updatedTimeEntry}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(&taskWithTimeEntry, nil).Once()
		mockTaskRepo.On("UpdateTimeEntry", mock.Anything, task.ID, timeEntryID.Hex(), mock.AnythingOfType("*models.TimeEntry")).
			Return(&taskAfterUpdate, nil).Once()

		// Create request body
		updatedTimeEntryJSON, _ := json.Marshal(updatedTimeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+task.ID+"/time-entries/"+timeEntryID.Hex(), bytes.NewBuffer(updatedTimeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.PUT("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.UpdateTimeEntry(c)
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
		assert.Equal(t, updatedStartDate, response.TimeEntries[0].StartDate)
		assert.Equal(t, updatedEndDate, response.TimeEntries[0].EndDate)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("task not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		taskID := primitive.NewObjectID().Hex()
		timeEntryID := primitive.NewObjectID()

		// Updated time entry values
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, taskID).Return(nil, nil).Once()

		// Create request body
		updatedTimeEntryJSON, _ := json.Marshal(updatedTimeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+taskID+"/time-entries/"+timeEntryID.Hex(), bytes.NewBuffer(updatedTimeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.PUT("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.UpdateTimeEntry(c)
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
		timeEntryID := primitive.NewObjectID()

		// Updated time entry values
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()

		// Create request body
		updatedTimeEntryJSON, _ := json.Marshal(updatedTimeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+task.ID+"/time-entries/"+timeEntryID.Hex(), bytes.NewBuffer(updatedTimeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.PUT("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.UpdateTimeEntry(c)
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
		timeEntryID := primitive.NewObjectID()

		// Updated time entry values
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()
		mockTaskRepo.On("UpdateTimeEntry", mock.Anything, task.ID, timeEntryID.Hex(), mock.AnythingOfType("*models.TimeEntry")).
			Return(nil, errors.New("no time entries found")).Once()

		// Create request body
		updatedTimeEntryJSON, _ := json.Marshal(updatedTimeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+task.ID+"/time-entries/"+timeEntryID.Hex(), bytes.NewBuffer(updatedTimeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.PUT("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.UpdateTimeEntry(c)
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
		timeEntryID := primitive.NewObjectID()

		// Updated time entry values
		updatedStartDate := time.Now().Add(30 * time.Minute).Format(time.RFC3339)
		updatedEndDate := time.Now().Add(2 * time.Hour).Format(time.RFC3339)
		updatedTimeEntry := &models.TimeEntry{
			ID:        &timeEntryID,
			StartDate: updatedStartDate,
			EndDate:   updatedEndDate,
		}

		// Mock repository calls
		mockTaskRepo.On("GetByID", mock.Anything, task.ID).Return(task, nil).Once()
		mockTaskRepo.On("UpdateTimeEntry", mock.Anything, task.ID, timeEntryID.Hex(), mock.AnythingOfType("*models.TimeEntry")).
			Return(nil, errors.New("database error")).Once()

		// Create request body
		updatedTimeEntryJSON, _ := json.Marshal(updatedTimeEntry)

		// Create test request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tasks/"+task.ID+"/time-entries/"+timeEntryID.Hex(), bytes.NewBuffer(updatedTimeEntryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Create a router with the route registered
		router := gin.New()
		controller := NewTaskController(mockTaskRepo, nil)
		router.PUT("/tasks/:id/time-entries/:entryId", func(c *gin.Context) {
			// Manually set the auth user on each request
			c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
			controller.UpdateTimeEntry(c)
		})

		// Call the endpoint
		router.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verify mocks were called
		mockTaskRepo.AssertExpectations(t)
	})
}
