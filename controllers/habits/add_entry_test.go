package habits

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAddHabitEntry(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful add habit entry", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		// Create test habit
		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = userID

		// Create test entry
		entryDate := primitive.NewDateTimeFromTime(time.Now())
		entry := models.HabitEntry{
			HabitID:   habitID,
			EntryDate: entryDate,
		}

		// Expected created entry
		createdEntry := models.HabitEntry{
			ID:        primitive.NewObjectID(),
			HabitID:   habitID,
			UserID:    userID,
			EntryDate: entryDate,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		// Mock the repository calls
		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil).Once()
		mockRepo.On("AddEntry", mock.Anything, mock.MatchedBy(func(e *models.HabitEntry) bool {
			// Match entry by HabitID
			return e.HabitID == habitID
		})).Return(&createdEntry, nil).Once()

		// Create request body
		entryJSON, _ := json.Marshal(entry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits/entry/add", bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.AddHabitEntry(ctx)

		// Assertions
		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.HabitEntry
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, createdEntry.ID, response.ID)
		assert.Equal(t, habitID, response.HabitID)
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, entryDate, response.EntryDate)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		habitID := primitive.NewObjectID()
		entry := models.HabitEntry{
			HabitID: habitID,
		}

		entryJSON, _ := json.Marshal(entry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits/entry/add", bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Call the endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits/entry/add", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.AddHabitEntry(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("habit not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		entry := models.HabitEntry{
			HabitID: habitID,
		}

		// Mock the repository call
		mockRepo.On("GetByID", mock.Anything, habitID).Return(nil, nil).Once()

		entryJSON, _ := json.Marshal(entry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits/entry/add", bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.AddHabitEntry(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		habitOwnerID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = habitOwnerID // Set a different user as owner

		entry := models.HabitEntry{
			HabitID: habitID,
		}

		// Mock the repository call
		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil).Once()

		entryJSON, _ := json.Marshal(entry)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits/entry/add", bytes.NewBuffer(entryJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.AddHabitEntry(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
