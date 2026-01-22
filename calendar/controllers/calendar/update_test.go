package calendar

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/calendar/tests/mocks"
	testmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateCalendarEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)
	created, _ := repoMock.Create(nil, &models.Calendar{ProdID: "orig", Version: "2.0"})
	created.ProdID = "updated"

	// Expect GetByID to be called and return the created calendar, and Update to succeed
	repoMock.On("GetByID", testmock.Anything, testmock.Anything).Return(created, nil)
	repoMock.On("Update", testmock.Anything, testmock.Anything, testmock.Anything).Return(created, nil)

	body := map[string]*models.Calendar{"calendar": created}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/calendar/"+created.ID.Hex(), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}
