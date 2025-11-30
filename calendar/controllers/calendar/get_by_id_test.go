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

func TestGetCalendarByIDEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	// create an entry directly in the mock by setting expectation to call through to the fallback
	cal := &models.Calendar{ProdID: "-//GET//EN", Version: "2.0"}
	// allow Create to fall through to internal store
	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)
	created, _ := repoMock.Create(nil, cal)

	// Expect GetByID to be called by the handler and return the created calendar
	repoMock.On("GetByID", testmock.Anything, testmock.Anything).Return(created, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/calendar/"+created.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
