package tags

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

func TestGetTagByID(t *testing.T) {
	_, mockRepo := setupTest()

	t.Run("successful get tag by id", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()
		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID // Set the tag owner

		mockRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, tag.Name, response.Name)
		assert.Equal(t, *tag.Color, *response.Color)
	})

	t.Run("tag not found", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		userID := primitive.NewObjectID()

		mockRepo.On("GetByID", mock.Anything, nonExistentID).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+nonExistentID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: nonExistentID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Tag not found", response["error"])
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		tagID := primitive.NewObjectID().Hex()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: tagID}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
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
		req, _ := http.NewRequest("GET", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("invalid tag ID format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidID := "not-a-valid-object-id"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+invalidID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		ctx.Params = []gin.Param{{Key: "id", Value: invalidID}}

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
