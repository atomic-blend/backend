package calendar

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/calendar/tests/mocks"
	testmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetAllCalendarsEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	// create some entries via fallback create
	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)
	repoMock.Create(nil, &models.Calendar{ProdID: "a", Version: "2.0"})
	repoMock.Create(nil, &models.Calendar{ProdID: "b", Version: "2.0"})

	// Expect GetAll to be called by the handler and return the created calendars
	repoMock.On("GetAll", testmock.Anything, testmock.Anything, int64(1), int64(10)).Return([]*models.Calendar{
		{ProdID: "a", Version: "2.0"},
		{ProdID: "b", Version: "2.0"},
	}, int64(2), nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/calendar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	repoMock.AssertExpectations(t)
}
