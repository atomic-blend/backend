package habits

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteHabit(t *testing.T) {
	_, mockRepo := setupTest() // Fixed to avoid unused variable

	t.Run("successful delete habit", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		// Ensure habit has a properly set UserID field
		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = userID // This needs to be a valid ObjectID

		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil)
		mockRepo.On("Delete", mock.Anything, habitID).Return(nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		habitID := primitive.NewObjectID().Hex()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/"+habitID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: habitID}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing habit ID", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		// No params set to simulate missing ID

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid habit ID format", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		invalidID := "invalid-id-format"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/"+invalidID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the controller directly with our context that has auth
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("habit not found", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, habitID).Return(nil, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
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
		req, _ := http.NewRequest("DELETE", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("database error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		habitID := primitive.NewObjectID()

		habit := createTestHabit()
		habit.ID = habitID
		habit.UserID = userID

		mockRepo.On("GetByID", mock.Anything, habitID).Return(habit, nil)
		mockRepo.On("Delete", mock.Anything, habitID).Return(errors.New("database error"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/"+habitID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: habitID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabit(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
