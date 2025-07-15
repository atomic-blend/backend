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

func TestTimeEntryController_Delete(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()
	entryID := primitive.NewObjectID()

	existingEntry := &models.TimeEntry{
		ID:        &entryID,
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
	}

	mockRepo.On("GetByID", mock.Anything, entryID.Hex()).Return(existingEntry, nil)
	mockRepo.On("Delete", mock.Anything, entryID.Hex()).Return(nil)

	router.DELETE("/time-entries/:id", func(c *gin.Context) {
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.Delete(c)
	})

	req, _ := http.NewRequest("DELETE", "/time-entries/"+entryID.Hex(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}
