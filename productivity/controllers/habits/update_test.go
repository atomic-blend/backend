package habits

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateHabit(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful update habit", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = userID

		updatedHabit := createTestHabit()
		updatedHabit.ID = habitID
		updatedHabit.UserID = userID
		newName := "Updated Habit"
		updatedHabit.Name = &newName
		newEmoji := "🎯"
		updatedHabit.Emoji = &newEmoji
		newFrequency := models.FrequencyWeekly
		updatedHabit.Frequency = &newFrequency

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Habit")).Return(updatedHabit, nil)

		habitJSON, _ := json.Marshal(updatedHabit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the handler directly
		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Habit
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *updatedHabit.Name, *response.Name)
		assert.Equal(t, *updatedHabit.Emoji, *response.Emoji)
		assert.Equal(t, userID, response.UserID)                      // Verify the habit owner hasn't changed
		assert.Equal(t, *updatedHabit.Frequency, *response.Frequency) // Verify frequency updated
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		habitID := primitive.NewObjectID().Hex()
		habit := createTestHabit()
		habitJSON, _ := json.Marshal(habit)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID, bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Call endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("habit not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		nonExistentID := primitive.NewObjectID()
		habit := createTestHabit()

		mockRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil)

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+nonExistentID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID.Hex()}}

		// Call the handler directly
		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		habitOwnerID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = habitOwnerID // Set a different user as owner

		updatedHabit := createTestHabit()
		updatedHabit.ID = habitID
		updatedHabit.UserID = wrongUserID // This should not matter as controller should check existing habit

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)

		habitJSON, _ := json.Marshal(updatedHabit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the handler directly
		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = userID

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the handler directly
		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid habit ID format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidID := "not-a-valid-object-id"
		habit := createTestHabit()

		habitJSON, _ := json.Marshal(habit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+invalidID, bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the handler directly
		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Validation test cases for updates
	t.Run("validation - update with missing required name", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = userID

		updatedHabit := createTestHabit()
		updatedHabit.ID = habitID
		updatedHabit.UserID = userID
		updatedHabit.Name = nil // Missing required name

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)

		habitJSON, _ := json.Marshal(updatedHabit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Name")
	})

	t.Run("validation - update with invalid frequency", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = userID

		updatedHabit := createTestHabit()
		updatedHabit.ID = habitID
		updatedHabit.UserID = userID
		invalidFreq := "annually" // Invalid frequency
		updatedHabit.Frequency = &invalidFreq

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)

		habitJSON, _ := json.Marshal(updatedHabit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Frequency")
		assert.Contains(t, response["error"], "validFrequency")
	})

	t.Run("validation - update with missing startDate", func(t *testing.T) {
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		existingHabit := createTestHabit()
		existingHabit.ID = habitID
		existingHabit.UserID = userID

		updatedHabit := createTestHabit()
		updatedHabit.ID = habitID
		updatedHabit.UserID = userID
		updatedHabit.StartDate = nil // Missing required startDate

		mockRepo.On("GetByID", mock.Anything, habitID).Return(existingHabit, nil)

		habitJSON, _ := json.Marshal(updatedHabit)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/habits/"+habitID.Hex(), bytes.NewBuffer(habitJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		controller := NewHabitController(mockRepo)
		controller.UpdateHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "StartDate")
	})
}
