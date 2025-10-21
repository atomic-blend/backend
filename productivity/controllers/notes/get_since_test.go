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

func TestGetNotesSince(t *testing.T) {
	_, mockNoteRepo := setupTest()

	t.Run("successful get notes since", func(t *testing.T) {
		userID := primitive.NewObjectID()
		notes := []*models.NoteEntity{
			createTestNote(),
			createTestNote(),
		}
		totalCount := int64(2)

		// Mock repository response
		mockNoteRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(notes, totalCount, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=2024-01-01T00:00:00Z&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, int64(1), response.Page)
		assert.Equal(t, int64(10), response.Size)
		assert.Equal(t, int64(1), response.TotalPages)
		assert.Len(t, response.Notes, 2)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("missing since parameter", func(t *testing.T) {
		userID := primitive.NewObjectID()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Missing required parameter: since", response["error"])
	})

	t.Run("invalid date format", func(t *testing.T) {
		userID := primitive.NewObjectID()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=invalid-date", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid date format")
	})

	t.Run("unauthorized request", func(t *testing.T) {
		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", controller.GetNotesSince)

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=2024-01-01T00:00:00Z", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Authentication required", response["error"])
	})

	t.Run("repository error", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Mock repository error
		mockNoteRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(nil, int64(0), assert.AnError).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=2024-01-01T00:00:00Z&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "assert.AnError")

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("successful get notes since with timezone offset", func(t *testing.T) {
		userID := primitive.NewObjectID()
		notes := []*models.NoteEntity{
			createTestNote(),
		}
		totalCount := int64(1)

		// Mock repository response
		mockNoteRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("*int64"), mock.AnythingOfType("*int64")).Return(notes, totalCount, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=2024-01-01T12:30:45%2B02:00&page=1&limit=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), response.TotalCount)
		assert.Len(t, response.Notes, 1)

		mockNoteRepo.AssertExpectations(t)
	})

	t.Run("successful get notes since without pagination", func(t *testing.T) {
		userID := primitive.NewObjectID()
		notes := []*models.NoteEntity{
			createTestNote(),
			createTestNote(),
		}
		totalCount := int64(2)

		// Mock repository response - when no pagination, page and limit are nil
		mockNoteRepo.On("GetSince", mock.Anything, userID, mock.AnythingOfType("time.Time"), (*int64)(nil), (*int64)(nil)).Return(notes, totalCount, nil).Once()

		controller := NewNoteController(mockNoteRepo)

		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/notes/since", func(c *gin.Context) {
			// Mock authentication
			authUser := &auth.UserAuthInfo{UserID: userID}
			c.Set("authUser", authUser)
			// No pagination parameters set
			controller.GetNotesSince(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/notes/since?since=2024-01-01T00:00:00Z", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedNoteResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), response.TotalCount)
		assert.Equal(t, int64(0), response.Page)
		assert.Equal(t, int64(0), response.Size)
		assert.Equal(t, int64(0), response.TotalPages)
		assert.Len(t, response.Notes, 2)

		mockNoteRepo.AssertExpectations(t)
	})
}
