package habits

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateHabit(t *testing.T) {
	_, mockRepo := setupTest()

	t.Run("successful create habit", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habit := createTestHabit()
		habit.UserID = userID // Should be overwritten by handler

		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Habit")).Return(habit, nil)

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.Habit
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *habit.Name, *response.Name)
		assert.Equal(t, userID, response.UserID) // Verify habit is owned by authenticated user
		assert.Equal(t, *habit.Frequency, *response.Frequency)
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		habit := createTestHabit()
		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Invalid JSON
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Validation test cases
	t.Run("validation - missing required name", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habit := createTestHabit()
		habit.Name = nil // Make name nil to trigger validation error

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Name")
	})

	t.Run("validation - missing required frequency", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habit := createTestHabit()
		habit.Frequency = nil // Make frequency nil to trigger validation error

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Frequency")
	})

	t.Run("validation - invalid frequency value", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habit := createTestHabit()
		invalidFreq := "yearly" // Not in ValidFrequencies list
		habit.Frequency = &invalidFreq

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Frequency")
		assert.Contains(t, response["error"], "validFrequency")
	})

	t.Run("validation - missing start date", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habit := createTestHabit()
		habit.StartDate = nil // Make start date nil to trigger validation error

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/habits", bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewHabitController(mockRepo)
		controller.CreateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "StartDate")
	})
}
