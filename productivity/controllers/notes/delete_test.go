package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeleteNote(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful delete note", func(t *testing.T) {
		userID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = userID
		noteID := existingNote.ID.Hex()

		// Mock repository responses
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()
		mockNoteRepo.On("Delete", mock.Anything, noteID).Return(nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note deleted successfully", response["message"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("delete note without authentication", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", controller.DeleteNote)

		req, _ := http.NewRequest(http.MethodDelete, "/notes/123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("delete note not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository response - note not found
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note not found", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("delete note forbidden - not owner", func(t *testing.T) {
		userID := primitive.NewObjectID()
		differentUserID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = differentUserID
		noteID := existingNote.ID.Hex()

		// Mock repository response
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "You don't have permission to delete this note", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("delete note repository error on GetByID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository error
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("delete note repository error on Delete", func(t *testing.T) {
		userID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = userID
		noteID := existingNote.ID.Hex()

		// Mock repository responses
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()
		mockNoteRepo.On("Delete", mock.Anything, noteID).Return(assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		req, _ := http.NewRequest(http.MethodDelete, "/notes/"+noteID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("delete note missing ID parameter", func(t *testing.T) {
		userID := primitive.NewObjectID()
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.DELETE("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.DeleteNote(c)
		})

		// Use an empty ID which should trigger the "Note ID is required" error
		req, _ := http.NewRequest(http.MethodDelete, "/notes/ ", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note ID is required", response["error"])
	})
}
