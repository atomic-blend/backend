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

func TestCreateCalendarEndpoint(t *testing.T) {
	repoMock := &mocks.MockCalendarRepository{}
	router := setupRouterWithMockRepoForTests(repoMock)

	// allow fallback behavior for Create
	repoMock.On("Create", testmock.Anything, testmock.Anything).Maybe().Return(nil, nil)

	cal := &models.Calendar{ProdID: "-//AB//EN", Version: "2.0"}
	body := map[string]*models.Calendar{"calendar": cal}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/calendar", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)
}
