package folder

import (
	"productivity/models"
	"productivity/tests/mocks"
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

func TestGetAllFolders(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create test data
		userID := primitive.NewObjectID()
		folders := []*models.Folder{
			createTestFolder(),
			createTestFolder(),
		}

		// Set the user ID for the test folders
		folders[0].UserID = userID
		folders[1].UserID = userID

		// Setup mock expectation
		mockRepo.On("GetAll", mock.Anything, userID).Return(folders, nil)

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/folders", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.GetAllFolders(c)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response body
		var responseFolders []*models.Folder
		err := json.Unmarshal(w.Body.Bytes(), &responseFolders)
		assert.NoError(t, err)
		assert.Equal(t, len(folders), len(responseFolders))

		mockRepo.AssertExpectations(t)
	})

	t.Run("No Auth", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without auth
		req, _ := http.NewRequest(http.MethodGet, "/folders", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the controller function
		controller.GetAllFolders(c)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		userID := primitive.NewObjectID()

		// Setup mock to return error
		mockRepo.On("GetAll", mock.Anything, userID).Return(nil, errors.New("database error"))

		// Create request
		req, _ := http.NewRequest(http.MethodGet, "/folders", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.GetAllFolders(c)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockRepo.AssertExpectations(t)
	})
}
