package userrole

import (
	"atomic-blend/backend/auth/models"
	"atomic-blend/backend/auth/tests/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetRoleByID(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockUserRoleRepository)
	controller := NewUserRoleController(mockRepo)
	router := gin.New()
	group := router.Group("/admin")
	controller.SetupRoutes(group)

	t.Run("successful get role by id", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()
		role := &models.UserRoleEntity{
			ID:   &roleID,
			Name: "test_role",
		}

		// Setup mock expectation
		mockRepo.On("GetByID", mock.Anything, roleID).Return(role, nil).Once()

		// Create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/user-roles/"+roleID.Hex(), nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.UserRoleEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, role.Name, response.Name)
		assert.Equal(t, role.ID.Hex(), response.ID.Hex())

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id format", func(t *testing.T) {
		// Create request with invalid ID
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/user-roles/invalid_id", nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid ID format", response["error"])
	})

	t.Run("role not found", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()

		// Setup mock expectation
		mockRepo.On("GetByID", mock.Anything, roleID).Return(nil, errors.New("user role not found")).Once()

		// Create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/user-roles/"+roleID.Hex(), nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User role not found", response["error"])

		// Verify mock
		mockRepo.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()

		// Setup mock expectation
		mockRepo.On("GetByID", mock.Anything, roleID).Return(nil, errors.New("internal error")).Once()

		// Create request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/admin/user-roles/"+roleID.Hex(), nil)

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "internal error", response["error"])

		// Verify mock
		mockRepo.AssertExpectations(t)
	})
}
