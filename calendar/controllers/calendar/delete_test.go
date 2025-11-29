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

func TestDeleteCalendarEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)
	created, _ := repoMock.Create(nil, &models.Calendar{ProdID: "to-delete", Version: "2.0"})

	// ensure GetByID returns the created calendar and Delete succeeds
	repoMock.On("GetByID", testmock.Anything, testmock.Anything).Return(created, nil)
	repoMock.On("Delete", testmock.Anything, testmock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/calendar/"+created.ID.Hex(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusNoContent, w.Code)
}
