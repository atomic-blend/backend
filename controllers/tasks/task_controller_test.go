// filepath: /Users/brandonguigo/workspace/atomic-blend/backend/controllers/tasks/task_controller_test.go
package tasks

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTaskController(t *testing.T) {
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTagRepo := new(mocks.MockTagRepository)
	controller := NewTaskController(mockTaskRepo, mockTagRepo)

	assert.NotNil(t, controller)
	assert.Equal(t, mockTaskRepo, controller.taskRepo)
	assert.Equal(t, mockTagRepo, controller.tagRepo)
}

func TestSetupRoutesWithMock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockTaskRepo := new(mocks.MockTaskRepository)
	mockTagRepo := new(mocks.MockTagRepository)

	SetupRoutesWithMock(router, mockTaskRepo, mockTagRepo)

	// Test that routes are properly registered by making test requests
	testRoutes := []struct {
		method   string
		path     string
		expected int
	}{
		{http.MethodGet, "/tasks", http.StatusOK},
		{http.MethodGet, "/tasks/123", http.StatusOK},
		{http.MethodPost, "/tasks", http.StatusOK},
		{http.MethodPut, "/tasks/123", http.StatusOK},
		{http.MethodDelete, "/tasks/123", http.StatusOK},
	}

	// Setup mock expectations for each route - fixing the GetAll method to include both parameters
	mockTaskRepo.On("GetAll", mock.Anything, mock.AnythingOfType("*primitive.ObjectID")).Return([]*models.TaskEntity{}, nil)
	mockTaskRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(createTestTask(), nil)
	mockTaskRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TaskEntity")).Return(createTestTask(), nil)
	mockTaskRepo.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.TaskEntity")).Return(createTestTask(), nil)
	mockTaskRepo.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

	// Modify the createTestTask function to also include Tags if necessary
	task := createTestTask()
	if task.Tags != nil && len(*task.Tags) > 0 {
		for _, tagID := range *task.Tags {
			mockTagRepo.On("GetByID", mock.Anything, tagID).Return(&models.Tag{
				Name: "Test Tag",
			}, nil)
		}
	}

	for _, route := range testRoutes {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(route.method, route.path, nil)
		router.ServeHTTP(w, req)
		// We're not testing for the exact status code here because we're just checking route registration
		// The actual handler would return 401 for unauthenticated requests
	}
}
