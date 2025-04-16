package tags

import (
	"atomic_blend_api/auth"
	"atomic_blend_api/models"
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
	_, mockTagRepo, mockTaskRepo := setupTest()

	t.Run("successful delete tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		// Ensure tag has a properly set UserID field
		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID // This needs to be a valid ObjectID

		// Create mock tasks with the tag
		task1 := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: &[]primitive.ObjectID{tagID, primitive.NewObjectID()},
		}
		task2 := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: &[]primitive.ObjectID{primitive.NewObjectID()},
		}
		tasks := []*models.TaskEntity{task1, task2}

		// Set up mock expectations
		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return *u == userID
		})).Return(tasks, nil).Once()

		// Task1 should be updated since it contains the tag
		mockTaskRepo.On("Update", mock.Anything, task1.ID, mock.MatchedBy(func(t *models.TaskEntity) bool {
			// The updated task should have the tag removed
			if t.Tags == nil {
				return false
			}
			for _, tag := range *t.Tags {
				if tag == tagID {
					return false
				}
			}
			return len(*t.Tags) == 1
		})).Return(task1, nil).Once()

		mockTagRepo.On("Delete", mock.Anything, tagID).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly with our context that has auth
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTagRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
	})

	t.Run("successful delete tag - some tasks don't have the tag", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		// Create tag
		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		// Create tasks with different tag scenarios
		task1 := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: &[]primitive.ObjectID{tagID, primitive.NewObjectID()}, // Has the tag being deleted
		}
		task2 := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: &[]primitive.ObjectID{primitive.NewObjectID()}, // Doesn't have the tag
		}
		task3 := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: nil, // No tags at all
		}
		tasks := []*models.TaskEntity{task1, task2, task3}

		// Set up mock expectations
		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return *u == userID
		})).Return(tasks, nil).Once()

		// Only task1 should be updated since it's the only one containing the tag
		mockTaskRepo.On("Update", mock.Anything, task1.ID, mock.MatchedBy(func(t *models.TaskEntity) bool {
			// The updated task should have the tag removed
			if t.Tags == nil {
				return false
			}
			for _, tag := range *t.Tags {
				if tag == tagID {
					return false
				}
			}
			return len(*t.Tags) == 1
		})).Return(task1, nil).Once()

		// Task2 and task3 should not be updated
		// No mock expectations for updating them

		mockTagRepo.On("Delete", mock.Anything, tagID).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
		mockTagRepo.AssertExpectations(t)
		mockTaskRepo.AssertExpectations(t)
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
		controller := NewTagController(mockTagRepo, mockTaskRepo)
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
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("tag not found - nil tag", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("tag not found - error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(nil, errors.New("tag not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
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

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: wrongUserID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("error getting tasks", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return *u == userID
		})).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("error updating task", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		task := &models.TaskEntity{
			ID:   primitive.NewObjectID().Hex(),
			User: userID,
			Tags: &[]primitive.ObjectID{tagID, primitive.NewObjectID()},
		}
		tasks := []*models.TaskEntity{task}

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return *u == userID
		})).Return(tasks, nil).Once()

		mockTaskRepo.On("Update", mock.Anything, task.ID, mock.Anything).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("internal server error on delete", func(t *testing.T) {
		userID := primitive.NewObjectID()
		tagID := primitive.NewObjectID()

		tag := createTestTag()
		tag.ID = &tagID
		tag.UserID = &userID

		// No tasks have this tag
		tasks := []*models.TaskEntity{}

		mockTagRepo.On("GetByID", mock.Anything, tagID).Return(tag, nil).Once()
		mockTaskRepo.On("GetAll", mock.Anything, mock.MatchedBy(func(u *primitive.ObjectID) bool {
			return *u == userID
		})).Return(tasks, nil).Once()
		mockTagRepo.On("Delete", mock.Anything, tagID).Return(errors.New("database error")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/tags/"+tagID.Hex(), nil)

		// Create a new context with the request and set auth user
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req
		ctx.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Copy params from original request to the new context
		ctx.Params = []gin.Param{{Key: "id", Value: tagID.Hex()}}

		// Call the controller directly
		controller := NewTagController(mockTagRepo, mockTaskRepo)
		controller.DeleteTag(ctx)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
