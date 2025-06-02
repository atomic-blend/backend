package timeentrycontroller

import (
	"atomic_blend_api/auth"
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

func TestTimeEntryController_GetAll(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()
	timeEntries := []*models.TimeEntry{
		{
			ID:        &primitive.ObjectID{},
			User:      &userID,
			StartDate: "2025-05-28T10:00:00Z",
			EndDate:   "2025-05-28T12:00:00Z",
		},
	}

	mockRepo.On("GetAll", mock.Anything, &userID).Return(timeEntries, nil)

	router.GET("/time-entries", func(c *gin.Context) {
		// Simulate auth middleware setting the user
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.GetAll(c)
	})

	req, _ := http.NewRequest("GET", "/time-entries", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestTimeEntryController_GetAll_Unauthorized(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	router.GET("/time-entries", controller.GetAll)

	req, _ := http.NewRequest("GET", "/time-entries", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
