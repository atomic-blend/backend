package tags

import (
	"productivity/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTagByID(t *testing.T) {
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful get tag by id", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		// Create test tag
		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, tag.Name, response.Name)
		assert.Equal(t, *tag.Color, *response.Color)
		assert.Equal(t, userID, *response.UserID)
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		tagID := primitive.NewObjectID().Hex()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID, nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Params = []gin.Param{{Key: "id", Value: tagID}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("missing tag id", func(t *testing.T) {
		userID := primitive.NewObjectID()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Empty tag ID param
		ctx.Params = []gin.Param{{Key: "id", Value: ""}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid tag id format", func(t *testing.T) {
		userID := primitive.NewObjectID()
		invalidTagID := "not-an-object-id"

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+invalidTagID, nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Set the invalid tag ID
		ctx.Params = []gin.Param{{Key: "id", Value: invalidTagID}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tag not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetTagByID(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
