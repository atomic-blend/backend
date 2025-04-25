package tags

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

func TestCreateTag(t *testing.T) {
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful create tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tag := createTestTag()
		tag.UserID = &userID // This should be overwritten by the handler

		mockTagRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Tag")).Return(tag, nil).Once()

		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, tag.Name, response.Name)
		assert.Equal(t, *tag.Color, *response.Color)
		assert.Equal(t, userID, *response.UserID) // Verify the tag is owned by the authenticated user
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		tag := createTestTag()
		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Invalid JSON
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required name field", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Create a tag without a name
		tag := createTestTag()
		tag.Name = ""
		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Field validation for 'Name' failed on the 'required' tag")
	})
}
