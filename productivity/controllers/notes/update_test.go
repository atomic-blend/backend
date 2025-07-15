package notes

import (
	"productivity/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateNote(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful update note", func(t *testing.T) {
		userID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = userID
		noteID := existingNote.ID.Hex()

		updatedNote := createTestNote()
		updatedNote.ID = existingNote.ID
		updatedNote.User = userID
		updatedNote.Title = stringPtr("Updated Title")
		updatedNote.Content = stringPtr("Updated Content")

		// Mock repository responses
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()
		mockNoteRepo.On("Update", mock.Anything, noteID, mock.AnythingOfType("*models.NoteEntity")).Return(updatedNote, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.NoteEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Title", *response.Title)
		assert.Equal(t, "Updated Content", *response.Content)
		assert.Equal(t, userID, response.User)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("update note without authentication", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", controller.UpdateNote)

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/123", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("update note not found", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository response - note not found
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Note not found", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("update note forbidden - not owner", func(t *testing.T) {
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
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "You don't have permission to update this note", response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("update note with invalid JSON", func(t *testing.T) {
		userID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = userID
		noteID := existingNote.ID.Hex()

		// Mock repository response
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "invalid character")

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("update note repository error on GetByID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		noteID := primitive.NewObjectID().Hex()

		// Mock repository error
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(jsonData))
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

	t.Run("update note repository error on Update", func(t *testing.T) {
		userID := primitive.NewObjectID()
		existingNote := createTestNote()
		existingNote.User = userID
		noteID := existingNote.ID.Hex()

		// Mock repository responses
		mockNoteRepo.On("GetByID", mock.Anything, noteID).Return(existingNote, nil).Once()
		mockNoteRepo.On("Update", mock.Anything, noteID, mock.AnythingOfType("*models.NoteEntity")).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.PUT("/notes/:id", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.UpdateNote(c)
		})

		requestNote := models.NoteEntity{
			Title:   stringPtr("Updated Title"),
			Content: stringPtr("Updated Content"),
		}

		jsonData, _ := json.Marshal(requestNote)
		req, _ := http.NewRequest(http.MethodPut, "/notes/"+noteID, bytes.NewBuffer(jsonData))
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
