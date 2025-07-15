package habits

import (
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteHabitEntry(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful delete habit entry", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		// Mock the repository calls
		mockRepo.On("DeleteEntry", mock.Anything, entryID).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/entry/delete/"+entryID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabitEntry(ctx)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "deleted successfully")
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		entryID := primitive.NewObjectID().Hex()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/entry/delete/"+entryID, nil)

		// Call the endpoint without authentication
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid entry ID format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidID := "not-a-valid-object-id"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/entry/delete/"+invalidID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabitEntry(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing entry ID", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/entry/delete/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		// No params set to simulate missing ID

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabitEntry(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("server error during delete", func(t *testing.T) {
		userID := primitive.NewObjectID()
		entryID := primitive.NewObjectID()

		// Mock the repository call to return an error
		mockRepo.On("DeleteEntry", mock.Anything, entryID).Return(assert.AnError).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/habits/entry/delete/"+entryID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: entryID.Hex()}}

		// Call the controller directly
		controller := NewHabitController(mockRepo)
		controller.DeleteHabitEntry(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
