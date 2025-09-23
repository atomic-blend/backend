package timeentrycontroller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"github.com/atomic-blend/backend/shared/middlewares/auth"
	"github.com/atomic-blend/backend/productivity/models"
	"github.com/atomic-blend/backend/productivity/tests/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTimeEntryController_Update(t *testing.T) {
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
		CreatedAt: "2025-05-28T10:00:00Z",
	}

	updateData := models.TimeEntry{
		StartDate: "2025-05-28T11:00:00Z",
		EndDate:   "2025-05-28T13:00:00Z",
	}

	updatedEntry := &models.TimeEntry{
		ID:        &entryID,
		User:      &userID,
		StartDate: "2025-05-28T11:00:00Z",
		EndDate:   "2025-05-28T13:00:00Z",
		CreatedAt: "2025-05-28T10:00:00Z",
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	mockRepo.On("GetByID", mock.Anything, entryID.Hex()).Return(existingEntry, nil)
	mockRepo.On("Update", mock.Anything, entryID.Hex(), mock.AnythingOfType("*models.TimeEntry")).Return(updatedEntry, nil)

	router.PUT("/time-entries/:id", func(c *gin.Context) {
		// TODO: replace that with grpc call
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.Update(c)
	})

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/time-entries/"+entryID.Hex(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}
