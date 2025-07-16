package timeentrycontroller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"atomic-blend/backend/productivity/auth"
	"atomic-blend/backend/productivity/models"
	"atomic-blend/backend/productivity/tests/mocks"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTimeEntryController_Create(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()
	entryID := primitive.NewObjectID()

	requestBody := models.TimeEntry{
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
	}

	expectedEntry := &models.TimeEntry{
		ID:        &entryID,
		User:      &userID,
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
		Timer:     &[]bool{true}[0],
		Pomodoro:  &[]bool{false}[0],
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(entry *models.TimeEntry) bool {
		return entry.User != nil && *entry.User == userID
	})).Return(expectedEntry, nil)

	router.POST("/time-entries", func(c *gin.Context) {
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.Create(c)
	})

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/time-entries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
	mockRepo.AssertExpectations(t)
}

func TestTimeEntryController_Create_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.MockTimeEntryRepository)
	controller := NewTimeEntryController(mockRepo)
	router := setupTestRouter()

	userID := primitive.NewObjectID()

	requestBody := models.TimeEntry{
		StartDate: "2025-05-28T10:00:00Z",
		EndDate:   "2025-05-28T12:00:00Z",
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.TimeEntry")).Return(nil, errors.New("database error"))

	router.POST("/time-entries", func(c *gin.Context) {
		c.Set("authUser", &auth.UserAuthInfo{UserID: userID})
		controller.Create(c)
	})

	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/time-entries", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	mockRepo.AssertExpectations(t)
}
