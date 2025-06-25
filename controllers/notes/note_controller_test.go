package notes

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewNoteController(t *testing.T) {
	mockNoteRepo := new(mocks.MockNoteRepository)
	controller := NewNoteController(mockNoteRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, mockNoteRepo, controller.noteRepo)
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockNoteRepo := new(mocks.MockNoteRepository)

	SetupRoutesWithMock(router, mockNoteRepo)

	// Test that routes are properly registered by making test requests
	testRoutes := []struct {
		method   string
		path     string
		expected int
	}{
		{http.MethodGet, "/notes", http.StatusOK},
		{http.MethodGet, "/notes/123", http.StatusOK},
		{http.MethodPost, "/notes", http.StatusOK},
		{http.MethodPut, "/notes/123", http.StatusOK},
		{http.MethodDelete, "/notes/123", http.StatusOK},
	}

	// Setup mock expectations for each route
	mockNoteRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*primitive.ObjectID")).Return([]*models.NoteEntity{}, nil)
	mockNoteRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(createTestNote(), nil)
	mockNoteRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.NoteEntity")).Return(createTestNote(), nil)
	mockNoteRepo.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.NoteEntity")).Return(createTestNote(), nil)
	mockNoteRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	for _, test := range testRoutes {
		req, _ := http.NewRequest(test.method, test.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// We expect 401 because no auth is provided, but this confirms the route exists
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	}
}

// Helper function to create a test note
func createTestNote() *models.NoteEntity {
	id := primitive.NewObjectID()
	title := "Test Note"
	content := "This is a test note"
	userID := primitive.NewObjectID()
	now := primitive.NewDateTimeFromTime(primitive.NewObjectID().Timestamp())

	return &models.NoteEntity{
		ID:        &id,
		Title:     &title,
		Content:   &content,
		User:      userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Helper function to set up test environment
func setupTest() (*gin.Context, *mocks.MockNoteRepository) {
	gin.SetMode(gin.TestMode)
	mockNoteRepo := new(mocks.MockNoteRepository)

	return nil, mockNoteRepo
}
