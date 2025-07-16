package userrole

import (
	"atomic-blend/backend/auth/models"
	"atomic-blend/backend/auth/tests/mocks"
	"bytes"
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

func TestUpdateRole(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockUserRoleRepository)
	controller := NewUserRoleController(mockRepo)
	router := gin.New()
	group := router.Group("/admin")
	controller.SetupRoutes(group)

	t.Run("successful update role", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()
		role := &models.UserRoleEntity{
			ID:   &roleID,
			Name: "updated_role",
		}

		// Setup mock expectations
		mockRepo.On("GetByID", mock.Anything, roleID).Return(role, nil)
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserRoleEntity")).Return(role, nil)

		// Create request
		roleJSON, _ := json.Marshal(role)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/admin/user-roles/"+roleID.Hex(), bytes.NewBuffer(roleJSON))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.UserRoleEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, role.Name, response.Name)
	})

	t.Run("invalid id format", func(t *testing.T) {
		// Create request with invalid ID
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/admin/user-roles/invalid_id", bytes.NewBuffer([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid ID format", response["error"])
	})

	t.Run("role not found", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()
		role := &models.UserRoleEntity{
			ID:   &roleID,
			Name: "updated_role",
		}

		// Setup mock expectations
		mockRepo.On("GetByID", mock.Anything, roleID).Return(nil, errors.New("user role not found"))

		// Create request
		roleJSON, _ := json.Marshal(role)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/admin/user-roles/"+roleID.Hex(), bytes.NewBuffer(roleJSON))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User role not found", response["error"])
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Setup
		roleID := primitive.NewObjectID()
		mockRepo.On("GetByID", mock.Anything, roleID).Return(&models.UserRoleEntity{ID: &roleID}, nil)

		// Create request with invalid JSON
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/admin/user-roles/"+roleID.Hex(), bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("update error", func(t *testing.T) {
		// Setup new mock for this test to ensure clean state
		mockRepo := new(mocks.MockUserRoleRepository)
		controller := NewUserRoleController(mockRepo)
		router := gin.New()
		group := router.Group("/admin")
		controller.SetupRoutes(group)

		roleID := primitive.NewObjectID()
		existingRole := &models.UserRoleEntity{
			ID:   &roleID,
			Name: "existing_role",
		}

		// Setup mock expectations with Once()
		mockRepo.On("GetByID", mock.Anything, roleID).Return(existingRole, nil).Once()
		mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.UserRoleEntity")).Return(nil, errors.New("update error")).Once()

		// Create request body
		updatedRole := &models.UserRoleEntity{
			ID:   &roleID,
			Name: "updated_role",
		}
		roleJSON, _ := json.Marshal(updatedRole)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/admin/user-roles/"+roleID.Hex(), bytes.NewBuffer(roleJSON))
		req.Header.Set("Content-Type", "application/json")

		// Perform request
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response gin.H
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "update error", response["error"])

		// Verify all mocks were called as expected
		mockRepo.AssertExpectations(t)
	})
}
