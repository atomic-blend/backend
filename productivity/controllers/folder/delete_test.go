package folder

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"productivity/tests/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteFolder(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Mock data
		userID := primitive.NewObjectID()
		folderID := primitive.NewObjectID()

		// Setup mock expectation
		mockRepo.On("Delete", mock.Anything, folderID).Return(nil)

		// Create request
		req, _ := http.NewRequest(http.MethodDelete, "/folders/"+folderID.Hex(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: folderID.Hex()}}

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.DeleteFolder(c)

		// Assert response
		assert.Equal(t, http.StatusNoContent, w.Code)

		mockRepo.AssertExpectations(t)
	})

	t.Run("No Auth", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without auth
		req, _ := http.NewRequest(http.MethodDelete, "/folders/123", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the controller function
		controller.DeleteFolder(c)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Folder ID", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without folder ID
		req, _ := http.NewRequest(http.MethodDelete, "/folders/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.DeleteFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Folder ID", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request with invalid folder ID
		req, _ := http.NewRequest(http.MethodDelete, "/folders/invalid-id", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "invalid-id"}}

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.DeleteFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		userID := primitive.NewObjectID()
		folderID := primitive.NewObjectID()

		// Setup mock to return error
		mockRepo.On("Delete", mock.Anything, folderID).Return(errors.New("database error"))

		// Create request
		req, _ := http.NewRequest(http.MethodDelete, "/folders/"+folderID.Hex(), nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: folderID.Hex()}}

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.DeleteFolder(c)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockRepo.AssertExpectations(t)
	})
}
