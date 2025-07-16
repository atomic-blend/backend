package folder

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	"atomic-blend/backend/productivity/models"
	"atomic-blend/backend/productivity/tests/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateFolder(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Mock input data
		userID := primitive.NewObjectID()
		folderID := primitive.NewObjectID()
		folderName := "Updated Test Folder"
		color := "#00FF00"

		inputFolder := models.Folder{
			Name:  folderName,
			Color: &color,
		}

		// Expected output with ID and timestamps
		now := primitive.NewDateTimeFromTime(time.Now())
		expectedFolder := models.Folder{
			ID:        &folderID,
			Name:      folderName,
			Color:     &color,
			UserID:    userID,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		// Setup mock expectation
		mockRepo.On("Update", mock.Anything, folderID, mock.AnythingOfType("*models.Folder")).Return(&expectedFolder, nil)

		// Create request
		jsonInput, _ := json.Marshal(inputFolder)
		req, _ := http.NewRequest(http.MethodPut, "/folders/"+folderID.Hex(), bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: folderID.Hex()}}

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		// Parse response body
		var responseFolder models.Folder
		err := json.Unmarshal(w.Body.Bytes(), &responseFolder)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolder.ID, responseFolder.ID)
		assert.Equal(t, expectedFolder.Name, responseFolder.Name)
		assert.Equal(t, expectedFolder.Color, responseFolder.Color)

		mockRepo.AssertExpectations(t)
	})

	t.Run("No Auth", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without auth
		req, _ := http.NewRequest(http.MethodPut, "/folders/123", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("No Folder ID", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without folder ID
		req, _ := http.NewRequest(http.MethodPut, "/folders/", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Folder ID", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request with invalid folder ID
		req, _ := http.NewRequest(http.MethodPut, "/folders/invalid-id", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: "invalid-id"}}

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		folderID := primitive.NewObjectID()

		// Invalid JSON
		req, _ := http.NewRequest(http.MethodPut, "/folders/"+folderID.Hex(), bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: folderID.Hex()}}

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		userID := primitive.NewObjectID()
		folderID := primitive.NewObjectID()
		folderName := "Updated Test Folder"

		inputFolder := models.Folder{
			Name: folderName,
		}

		// Setup mock to return error
		mockRepo.On("Update", mock.Anything, folderID, mock.AnythingOfType("*models.Folder")).
			Return(nil, errors.New("database error"))

		// Create request
		jsonInput, _ := json.Marshal(inputFolder)
		req, _ := http.NewRequest(http.MethodPut, "/folders/"+folderID.Hex(), bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = []gin.Param{{Key: "id", Value: folderID.Hex()}}

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.UpdateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockRepo.AssertExpectations(t)
	})
}
