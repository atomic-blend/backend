package habits

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

func TestGetHabitByID(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful get habit by id", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()
		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = userID // Set the habit owner

		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		// Copy headers and params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.GetHabitByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Habit
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *habit.Name, *response.Name)
		assert.Equal(t, *habit.Emoji, *response.Emoji)
		assert.Equal(t, *habit.Frequency, *response.Frequency)
		assert.Equal(t, habit.Reminders, response.Reminders)
	})

	t.Run("habit not found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		userID := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits/"+nonExistentID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.GetHabitByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Habit not found", response["error"])
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		habitID := primitive.NewObjectID().Hex()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits/"+habitID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		habitOwnerID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = habitOwnerID // Set a different user as owner

		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		// Copy headers and params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.GetHabitByID(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid habit ID format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidID := "not-a-valid-object-id"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/habits/"+invalidID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.GetHabitByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid habit ID format", response["error"])
	})
}
