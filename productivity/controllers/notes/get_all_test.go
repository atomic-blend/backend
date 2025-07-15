package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"productivity/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllNotes(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful get all notes", func(t *testing.T) {
		userID := primitive.NewObjectID()
		notes := []*models.NoteEntity{
			createTestNote(),
			createTestNote(),
		}

		// Mock repository response
		mockNoteRepo.On("GetAll", mock.Anything, &userID).Return(notes, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetAllNotes(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.NoteEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get all notes without authentication", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", controller.GetAllNotes)

		req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("get all notes repository error", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Mock repository error
		mockNoteRepo.On("GetAll", mock.Anything, &userID).Return(nil, assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetAllNotes(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, assert.AnError.Error(), response["error"])

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get all notes empty result", func(t *testing.T) {
		userID := primitive.NewObjectID()
		var emptyNotes []*models.NoteEntity

		// Mock repository response with empty result
		mockNoteRepo.On("GetAll", mock.Anything, &userID).Return(emptyNotes, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetAllNotes(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*models.NoteEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 0)

		mockNoteRepo.AssertExpectations(t)
	})
}
