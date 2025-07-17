package notes

import (
	"bytes"
	"encoding/json"
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

func TestCreateNote(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful create note", func(t *testing.T) {
		// Create authenticated user
		userID := primitive.NewObjectID()
		note := createTestNote()

		// Set note user to match the authenticated user
		note.User = userID

		// Mock repository response
		mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NoteEntity")).Return(note, nil).Once()

		// Create controller
		controller := NewNoteController(mockNoteRepo)

		// Setup request
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.CreateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   note.Title,
			Content: note.Content,
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.NoteEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *note.Title, *response.Title)
		assert.Equal(t, *note.Content, *response.Content)
		assert.Equal(t, userID, response.User)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("create note without authentication", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/notes", controller.CreateNote)

		requestNote := models.NoteEntity{
			Title:   stringPtr("Test Note"),
			Content: stringPtr("Test Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("create note with invalid JSON", func(t *testing.T) {
		userID := primitive.NewObjectID()
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.CreateNote(c)
		})

		req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid character")
	})

	t.Run("create note repository error", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Mock repository error
		mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NoteEntity")).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.POST("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.CreateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Test Note"),
			Content: stringPtr("Test Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPost, "/notes", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response["error"])

		mockNoteRepo.AssertExpectations(t)
	})
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
