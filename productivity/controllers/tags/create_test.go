package tags

import (
	"productivity/models"
	"productivity/tests/utils/inmemorymongo"
	"productivity/utils/db"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTag(t *testing.T) {
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful create tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tag := createTestTag()
		tag.UserID = &userID // This should be overwritten by the handler

		// Mock GetAll to return fewer than 5 tags (no subscription needed)
		mockTagRepo.On("GetAll", mock.Anything, &userID).Return([]*models.Tag{}, nil).Once()
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

	t.Run("forbidden - user has 5 tags and is not subscribed", func(t *testing.T) {
		// Setup in-memory MongoDB for subscription check
		mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
		require.NoError(t, err)
		defer mongoServer.Stop()

		client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
		require.NoError(t, err)
		defer client.Disconnect(context.Background())

		// Set global database for subscription function
		originalDB := db.Database
		db.Database = client.Database("test_db")
		defer func() { db.Database = originalDB }()

		userID := primitive.NewObjectID()
		tag := createTestTag()

		// Create 5 existing tags to simulate user at limit
		existingTags := make([]*models.Tag, 5)
		for i := 0; i < 5; i++ {
			existingTags[i] = createTestTag()
			existingTags[i].UserID = &userID
		}

		// Mock GetAll to return 5 tags
		mockTagRepo.On("GetAll", mock.Anything, &userID).Return(existingTags, nil).Once()

		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "You must be subscribed to create more than 5 tags", response["error"])
	})

	t.Run("error when GetAll fails", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tag := createTestTag()

		// Mock GetAll to return an error (this happens before subscription check)
		mockTagRepo.On("GetAll", mock.Anything, &userID).Return(nil, assert.AnError).Once()

		tagJSON, _ := json.Marshal(tag)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/tags", bytes.NewBuffer(tagJSON))
		req.Header.Set("Content-Type", "application/json")

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.CreateTag(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "assert.AnError")
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

		// Mock GetAll to return fewer than 5 tags (no subscription needed)
		mockTagRepo.On("GetAll", mock.Anything, &userID).Return([]*models.Tag{}, nil).Once()

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

		// Mock GetAll to return fewer than 5 tags (no subscription needed)
		mockTagRepo.On("GetAll", mock.Anything, &userID).Return([]*models.Tag{}, nil).Once()

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
