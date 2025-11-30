package calendar

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/calendar/tests/mocks"
	testmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetSinceEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)
	created, _ := repoMock.Create(nil, &models.Calendar{ProdID: "since-test", Version: "2.0"})

	// Expect GetSince to be called and return the created calendar
	repoMock.On("GetSince", testmock.Anything, testmock.Anything, testmock.Anything, int64(1), int64(10)).Return([]*models.Calendar{created}, int64(1), nil)

	since := time.Now().Add(-time.Minute).UTC().Format(time.RFC3339)
	req := httptest.NewRequest(http.MethodGet, "/calendar/since?since="+since, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// quick decode to ensure returned payload includes created item
	var resp struct {
		Calendars  []*models.Calendar `json:"calendars"`
		TotalCount int64              `json:"total_count"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.GreaterOrEqual(t, resp.TotalCount, int64(1))

	// ensure the created ID is present
	found := false
	for _, c := range resp.Calendars {
		if c.ID != nil && *c.ID == *created.ID {
			found = true
			break
		}
	}
	require.True(t, found)
}
