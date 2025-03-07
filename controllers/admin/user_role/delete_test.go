package user_role

import (
	"atomic_blend_api/tests/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupDeleteTest() (*gin.Engine, *mocks.MockUserRoleRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockUserRoleRepository)
	roleController := NewUserRoleController(mockRepo)

	adminRoutes := router.Group("/admin/user-roles")
	{
		adminRoutes.DELETE("/:id", roleController.DeleteRole)
	}

	return router, mockRepo
}

func TestDeleteRole(t *testing.T) {
	tests := []struct {
		name       string
		roleID     string
		mockError  error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Invalid ID format",
			roleID:     "invalid-id",
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"Invalid ID format"}`,
		},
		{
			name:       "Role not found",
			roleID:     primitive.NewObjectID().Hex(),
			mockError:  errors.New("user role not found"),
			wantStatus: http.StatusNotFound,
			wantBody:   `{"error":"User role not found"}`,
		},
		{
			name:       "Successful deletion",
			roleID:     primitive.NewObjectID().Hex(),
			mockError:  nil,
			wantStatus: http.StatusNoContent,
			wantBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockRepo := setupDeleteTest()

			if tt.roleID != "invalid-id" {
				objID, _ := primitive.ObjectIDFromHex(tt.roleID)
				mockRepo.On("Delete", mock.Anything, objID).Return(tt.mockError)
			}

			// Create request
			req := httptest.NewRequest(http.MethodDelete, "/admin/user-roles/"+tt.roleID, nil)
			w := httptest.NewRecorder()

			// Perform request
			router.ServeHTTP(w, req)

			// Assert results
			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.wantBody != "" {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
