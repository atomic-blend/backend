package habits

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEditHabitEntry(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful edit habit entry", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		// Create test habit
		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = userID

		// Create test entry with updated data
		newEntryDate := primitive.NewDateTimeFromTime(time.Now().Add(24 * time.Hour)) // 1 day later
		updatedEntry := models.HabitEntry{
			ID:        entryID,
			HabitID:   habitID,
			EntryDate: newEntryDate,
		}

		// Result entry after update
		resultEntry := models.HabitEntry{
			ID:        entryID,
			HabitID:   habitID,
			UserID:    userID,
			EntryDate: newEntryDate,
			CreatedAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339), // Created 1 hour ago
			UpdatedAt: time.Now().Format(time.RFC3339),                     // Updated now
		}

		// Mock the repository calls
		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil).Once()
		mockRepo.On("UpdateEntry", mock.Anything, mock.MatchedBy(func(e *models.HabitEntry) bool {
			// Match entry by ID and HabitID
			return e.ID == entryID && e.HabitID == habitID
		})).Return(&resultEntry, nil).Once()

		// Create request body
		entryJSON, _ := json.Marshal(updatedEntry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+entryID.Hex(), bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.EditHabitEntry(ctx)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.HabitEntry
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, entryID, response.ID)
		assert.Equal(t, habitID, response.HabitID)
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, newEntryDate, response.EntryDate)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		entryID := primitive.NewObjectID().Hex()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+entryID, nil)

		// Call the endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid entry ID format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidID := "not-a-valid-object-id"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+invalidID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.EditHabitEntry(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+entryID.Hex(), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.EditHabitEntry(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("habit not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		updatedEntry := models.HabitEntry{
			ID:      entryID,
			HabitID: habitID,
		}

		// Mock the repository call
		mockRepo.On("GetByID", mock.Anything, habitID).Return(nil, nil).Once()

		entryJSON, _ := json.Marshal(updatedEntry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+entryID.Hex(), bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.EditHabitEntry(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		habitOwnerID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = habitOwnerID // Set a different user as owner

		updatedEntry := models.HabitEntry{
			ID:      entryID,
			HabitID: habitID,
		}

		// Mock the repository call
		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil).Once()

		entryJSON, _ := json.Marshal(updatedEntry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/entry/edit/"+entryID.Hex(), bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.EditHabitEntry(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
