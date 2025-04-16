package tags

import (
	"atomic_blend_api/auth"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteTag(t *testing.T) {
	_, mockRepo := setupTest()

	t.Run("successful delete tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		// Ensure tag has a properly set UserID field
		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID // This needs to be a valid ObjectID

		mockRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockRepo.On("Delete", mock.Anything, tagID).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		tagID := primitive.NewObjectID().Hex()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: tagID}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing tag ID", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Empty tag ID
		ctx.Params = []gin.Param{{Key: "id", Value: ""}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tag not found - nil tag", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, tagID).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("tag not found - error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, tagID).Return(nil, errors.New("tag not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		tagOwnerID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &tagOwnerID // Set a different user as owner

		mockRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("internal server error on delete", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		mockRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockRepo.On("Delete", mock.Anything, tagID).Return(errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
