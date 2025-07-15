package timeentrycontroller

import (
	"net/http"
	"net/http/httptest"
	"productivity/auth"
	"productivity/models"
	"productivity/tests/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTimeEntryController_GetByID(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()
	entryID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		ID:        &entryID,
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
	}

	mockRepo.On("GetByID", mock.Anything, entryID.Hex()).Return(timeEntry, nil)

	router.GET("/time-entries/:id", func(c *gin.Context) {
		// TODO: replace that with grpc call
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.GetByID(c)
	})

	req, _ := http.NewRequest("GET", "/time-entries/"+entryID.Hex(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestTimeEntryController_GetByID_Forbidden(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()
	otherUserID := primitive.NewObjectID()
	entryID := primitive.NewObjectID()
	timeEntry := &models.TimeEntry{
		ID:        &entryID,
		User:      &otherUserID, // Different user owns this entry
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
	}

	mockRepo.On("GetByID", mock.Anything, entryID.Hex()).Return(timeEntry, nil)

	router.GET("/time-entries/:id", func(c *gin.Context) {
		//TODO: replace that with grpc call
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.GetByID(c)
	})

	req, _ := http.NewRequest("GET", "/time-entries/"+entryID.Hex(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusForbidden, resp.Code)
	mockRepo.AssertExpectations(t)
}
