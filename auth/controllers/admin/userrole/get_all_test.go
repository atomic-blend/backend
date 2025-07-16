package userrole

import (
	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/tests/mocks"
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

func setupGetAllTest() (*gin.Engine, *mocks.MockUserRoleRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockRepo := new(mocks.MockUserRoleRepository)
	roleController := NewUserRoleController(mockRepo)

	adminRoutes := router.Group("/admin/user-roles")
	{
		adminRoutes.GET("", roleController.GetAllRoles)
	}

	return router, mockRepo
}

func TestGetAllRoles(t *testing.T) {
	// Create fixed ObjectIDs for predictable JSON comparison
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()

	tests := []struct {
		name       string
		mockRoles  []*models.UserRoleEntity // Change to pointer slice
		mockError  error
		wantStatus int
		wantBody   string
	}{
		{
			name: "Successful retrieval",
			mockRoles: []*models.UserRoleEntity{ // Update test data to use pointers
				{
					ID:   &id1,
					Name: "Admin",
				},
				{
					ID:   &id2,
					Name: "User",
				},
			},
			mockError:  nil,
			wantStatus: http.StatusOK,
			wantBody:   "", // We'll compare the actual JSON in the test
		},
		{
			name:       "Database error",
			mockRoles:  nil,
			mockError:  errors.New("database error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"database error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockRepo := setupGetAllTest()

			mockRepo.On("GetAll", mock.Anything).Return(tt.mockRoles, tt.mockError)

			req := httptest.NewRequest(http.MethodGet, "/admin/user-roles", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if tt.mockError != nil {
				assert.JSONEq(t, tt.wantBody, w.Body.String())
			} else {
				var response []*models.UserRoleEntity
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockRoles), len(response))
				for i, role := range response {
					assert.Equal(t, tt.mockRoles[i].ID, role.ID)
					assert.Equal(t, tt.mockRoles[i].Name, role.Name)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
