package userrole

import (
	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestRole() *models.UserRoleEntity {
	return &models.UserRoleEntity{
		Name: "Test Role",
	}
}

func setupTest() (*gin.Engine, *mocks.MockUserRoleRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockUserRoleRepository)
	roleController := NewUserRoleController(mockRepo)

	adminRoutes := router.Group("/admin/user-roles")
	{
		adminRoutes.POST("", roleController.CreateRole)
	}

	return router, mockRepo
}

func TestCreateRole(t *testing.T) {
	router, mockRepo := setupTest()

	t.Run("successful role creation", func(t *testing.T) {
		role := createTestRole()
		mockRepo.On("GetByName", mock.Anything, role.Name).Return(nil, nil)
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.UserRoleEntity")).Return(role, nil)

		roleJSON, _ := json.Marshal(role)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/admin/user-roles", bytes.NewBuffer(roleJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response models.UserRoleEntity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, role.Name, response.Name)
	})

	t.Run("duplicate role name", func(t *testing.T) {
		role := createTestRole()
		existingRole := createTestRole()

		// Reset mock expectations
		mockRepo := new(mocks.MockUserRoleRepository)
		roleController := NewUserRoleController(mockRepo)

		// Setup router with new controller
		router := gin.New()
		router.POST("/admin/user-roles", roleController.CreateRole)

		// Mock GetByName to return an existing role (indicating duplicate)
		mockRepo.On("GetByName", mock.Anything, role.Name).Return(existingRole, nil)

		roleJSON, _ := json.Marshal(role)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/admin/user-roles", bytes.NewBuffer(roleJSON))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "already exists")
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/admin/user-roles", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
