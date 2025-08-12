package folder

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/tests/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateFolder(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Mock input data
		userID := primitive.NewObjectID()
		folderName := "New Test Folder"
		color := "#00FF00"

		inputFolder := models.Folder{
			Name:  folderName,
			Color: &color,
		}

		// Expected output with ID and timestamps
		id := primitive.NewObjectID()
		now := primitive.NewDateTimeFromTime(time.Now())
		expectedFolder := models.Folder{
			ID:        &id,
			Name:      folderName,
			Color:     &color,
			UserID:    userID,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		// Setup mock expectation
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Folder")).Return(&expectedFolder, nil)

		// Create request
		jsonInput, _ := json.Marshal(inputFolder)
		req, _ := http.NewRequest(http.MethodPost, "/folders", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.CreateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		// Parse response body
		var responseFolder models.Folder
		err := json.Unmarshal(w.Body.Bytes(), &responseFolder)
		assert.NoError(t, err)
		assert.Equal(t, expectedFolder.ID, responseFolder.ID)
		assert.Equal(t, expectedFolder.Name, responseFolder.Name)
		assert.Equal(t, expectedFolder.Color, responseFolder.Color)
		assert.Equal(t, expectedFolder.UserID, responseFolder.UserID)

		mockRepo.AssertExpectations(t)
	})

	t.Run("No Auth", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Create request without auth
		req, _ := http.NewRequest(http.MethodPost, "/folders", nil)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Call the controller function
		controller.CreateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Input", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Invalid JSON
		req, _ := http.NewRequest(http.MethodPost, "/folders", bytes.NewBufferString("{invalid json}"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		userID := primitive.NewObjectID()
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.CreateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockRepo := new(mocks.MockFolderRepository)
		controller := NewFolderController(mockRepo)

		// Mock input data
		userID := primitive.NewObjectID()
		folderName := "New Test Folder"

		inputFolder := models.Folder{
			Name: folderName,
		}

		// Setup mock to return error
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Folder")).
			Return(nil, errors.New("database error"))

		// Create request
		jsonInput, _ := json.Marshal(inputFolder)
		req, _ := http.NewRequest(http.MethodPost, "/folders", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		// Set mock auth user in context
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})

		// Call the controller function
		controller.CreateFolder(c)

		// Assert response
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockRepo.AssertExpectations(t)
	})
}
