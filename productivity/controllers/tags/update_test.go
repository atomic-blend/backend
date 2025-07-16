package tags

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/atomic-blend/backend/productivity/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateTag(t *testing.T) {
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful update tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		existingTag := createTestTag()
		existingTag.ID = &tagID
		existingTag.UserID = &userID

		updatedTag := createTestTag()
		updatedTag.ID = &tagID
		updatedTag.UserID = &userID
		updatedTag.Name = "Updated Tag"
		newColor := "#00FF00"
		updatedTag.Color = &newColor

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(existingTag, nil).Once()
		mockTagRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Tag")).Return(updatedTag, nil).Once()

		tagJSON, _ := json.Marshal(updatedTag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+tagID.Hex(), bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the handler directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updatedTag.Name, response.Name)
		assert.Equal(t, *updatedTag.Color, *response.Color)
		assert.Equal(t, userID, *response.UserID) // Verify the tag owner hasn't changed
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		tagID := primitive.NewObjectID().Hex()
		tag := createTestTag()
		tagJSON, _ := json.Marshal(tag)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+tagID, bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: tagID}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("tag not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		nonExistentID := primitive.NewObjectID()
		tag := createTestTag()

		mockTagRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil).Once()

		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+nonExistentID.Hex(), bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID.Hex()}}

		// Call the handler directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("forbidden access - wrong user", func(t *testing.T) {
		wrongUserID := primitive.NewObjectID()
		tagOwnerID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		existingTag := createTestTag()
		existingTag.ID = &tagID
		existingTag.UserID = &tagOwnerID // Set a different user as owner

		updatedTag := createTestTag()
		updatedTag.ID = &tagID
		updatedTag.Name = "Updated Tag"

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(existingTag, nil).Once()

		tagJSON, _ := json.Marshal(updatedTag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+tagID.Hex(), bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the handler directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		existingTag := createTestTag()
		existingTag.ID = &tagID
		existingTag.UserID = &userID

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(existingTag, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+tagID.Hex(), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the handler directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required name field", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		existingTag := createTestTag()
		existingTag.ID = &tagID
		existingTag.UserID = &userID

		updatedTag := createTestTag()
		updatedTag.ID = &tagID
		updatedTag.UserID = &userID
		updatedTag.Name = "" // Missing required name

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(existingTag, nil).Once()

		tagJSON, _ := json.Marshal(updatedTag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/tags/"+tagID.Hex(), bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a context and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the handler directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.UpdateTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
