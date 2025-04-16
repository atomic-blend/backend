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

func TestGetAllTags(t *testing.T) {
	_, mockRepo := setupTest()

	t.Run("successful get all tags - with tags", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Create test tags
		tag1 := createTestTag()
		tag1ID := primitive.NewObjectID()
		tag1.ID = &tag1ID
		tag1.UserID = &userID

		tag2 := createTestTag()
		tag2ID := primitive.NewObjectID()
		tag2.ID = &tag2ID
		tag2.UserID = &userID

		tags := []*models.Tag{tag1, tag2}

		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(tags, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(response))
	})

	t.Run("successful get all tags - empty list", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Return nil to simulate no tags
		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		// Check that the response is an empty array, not null
		assert.Equal(t, "[]", w.Body.String())
	})

	t.Run("unauthorized access - no auth user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Simulate database error
		mockRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(userID *primitive.ObjectID) bool {
			return true
		})).Return(nil, assert.AnError).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTagController(mockRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
