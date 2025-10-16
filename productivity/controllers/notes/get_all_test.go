package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/shared/middlewares/auth"

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
		totalCount := int64(2)

		// Mock repository response
		mockNoteRepo.On("GetAll", mock.Anything, &userID, mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(notes, totalCount, nil).Once()

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

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Notes, 2)
		assert.Equal(t, totalCount, response.TotalCount)

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
		mockNoteRepo.On("GetAll", mock.Anything, &userID, mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(nil, int64(0), assert.AnError).Once()

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
		totalCount := int64(0)

		// Mock repository response with empty result
		mockNoteRepo.On("GetAll", mock.Anything, &userID, mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(emptyNotes, totalCount, nil).Once()

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

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Notes, 0)
		assert.Equal(t, totalCount, response.TotalCount)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get all notes with pagination", func(t *testing.T) {
		userID := primitive.NewObjectID()
		notes := []*models.NoteEntity{
			createTestNote(),
			createTestNote(),
		}
		totalCount := int64(10) // Total count is higher than returned notes

		// Mock repository response with pagination
		mockNoteRepo.On("GetAll", mock.Anything, &userID, mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(notes, totalCount, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetAllNotes(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes?page=1&limit=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Notes, 2)
		assert.Equal(t, totalCount, response.TotalCount)
		assert.Equal(t, int64(1), response.Page)
		assert.Equal(t, int64(2), response.Size)
		assert.Equal(t, int64(5), response.TotalPages) // 10 total / 2 per page = 5 pages

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("get all notes with invalid pagination parameters", func(t *testing.T) {
		userID := primitive.NewObjectID()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetAllNotes(c)
		})

		// Test with invalid page parameter
		req, _ := http.NewRequest(http.MethodGet, "/notes?page=0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid page parameter", response["error"])

		// Test with invalid limit parameter
		req, _ = http.NewRequest(http.MethodGet, "/notes?limit=0", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid limit parameter", response["error"])
	})
}
