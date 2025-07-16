package folder

import (
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

func TestNewFolderController(t *testing.T) {
	mockRepo := new(mocks.MockFolderRepository)
	controller := NewFolderController(mockRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, mockRepo, controller.folderRepo)
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockFolderRepository)

	SetupRoutesWithMock(router, mockRepo)

	// Test that routes are properly registered by making test requests
	testRoutes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/folders"},
		{http.MethodPost, "/folders"},
		{http.MethodPut, "/folders/123"},
		{http.MethodDelete, "/folders/123"},
	}

	// Mock authentication middleware
	router.Use(func(c *gin.Context) {
		// Set mock auth user in context
		userID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		c.Next()
	})

	// Setup mock expectations for each route
	mockRepo.On("GetAll", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return([]*models.Folder{}, nil)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Folder")).Return(createTestFolder(), nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("primitive.ObjectID"), mock.AnythingOfType("*models.Folder")).Return(createTestFolder(), nil)
	mockRepo.On("Delete", mock.Anything, mock.AnythingOfType("primitive.ObjectID")).Return(nil)

	for _, route := range testRoutes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		router.ServeHTTP(w, req)

		// We don't expect 404s which would indicate the route isn't registered
		assert.NotEqual(t, http.StatusNotFound, w.Code, "Route not found: %s %s", route.method, route.path)
	}
}

// Helper function to create a test folder for test data
func createTestFolder() *models.Folder {
	id := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	name := "Test Folder"
	color := "#FF0000"
	now := primitive.NewDateTimeFromTime(time.Now())

	return &models.Folder{
		ID:        &id,
		Name:      name,
		Color:     &color,
		UserID:    userID,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
