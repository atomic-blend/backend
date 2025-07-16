package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	"atomic-blend/backend/productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetNoteByID(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful get note by ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		note := createTestNote()
		note.User = userID
		noteID := note.ID.Hex()

		// Mock repository response
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(note, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNoteByID(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.NoteEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, *note.Title, *response.Title)
		assert.Equal(t, *note.Content, *response.Content)
		assert.Equal(t, userID, response.User)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get note by ID without authentication", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", controller.GetNoteByID)

		req, _ := http.NewRequest(http.MethodGet, "/notes/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("get note by ID not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository response - note not found
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNoteByID(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note not found", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get note by ID forbidden - not owner", func(t *testing.T) {
		userID := primitive.NewObjectID()
		differentUserID := primitive.NewObjectID()
		note := createTestNote()
		note.User = differentUserID
		noteID := note.ID.Hex()

		// Mock repository response
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(note, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNoteByID(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "You don't have permission to access this note", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get note by ID repository error", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository error
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNoteByID(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get note by ID missing ID parameter", func(t *testing.T) {
		userID := primitive.NewObjectID()
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNoteByID(c)
		})

		// Use an empty ID which should trigger the "Note ID is required" error
		req, _ := http.NewRequest(http.MethodGet, "/notes/ ", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note ID is required", response["error"])
	})
}
