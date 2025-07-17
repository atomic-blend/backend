package tags

import (
	"encoding/json"
	"errors"
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

func TestGetAllTags(t *testing.T) {
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful get all tags", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		// Generate a few tags
		tag1 := createTestTag()
		tag1.Name = "Tag 1"
		id1 := primitive.NewObjectID()
		tag1.ID = &id1
		tag1.UserID = &userID

		tag2 := createTestTag()
		tag2.Name = "Tag 2"
		id2 := primitive.NewObjectID()
		tag2.ID = &id2
		tag2.UserID = &userID

		tags := []*models.Tag{tag1, tag2}

		mockTagRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return u != nil && *u == userID
		})).Return(tags, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		assert.Equal(t, "Tag 1", response[0].Name)
		assert.Equal(t, "Tag 2", response[1].Name)
	})

	t.Run("unauthorized - no auth user", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request but no auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		mockTagRepo.On("GetAll", mock.Anything, mock.Anything).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("no content - empty tags", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()

		mockTagRepo.On("GetAll", mock.Anything, mock.Anything).Return([]*models.Tag{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tags", nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.GetAllTags(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []*models.Tag
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)
	})
}
