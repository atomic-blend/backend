package repositories

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/atomic-blend/backend/calendar/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupCalendarTest(t *testing.T) (CalendarRepositoryInterface, func()) {
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	db := client.Database("test_db")
	repo := NewCalendarRepository(db)

	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestCalendar() *models.Calendar {
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Calendar{
		ProdID:    "-//Atomic Blend//EN",
		Version:   "2.0",
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func TestCalendarRepository_CreateGetUpdateDelete(t *testing.T) {
	repo, cleanup := setupCalendarTest(t)
	defer cleanup()

	userID := primitive.NewObjectID()

	cal := createTestCalendar()
	cal.UserID = &userID

	// use a gin context for repository calls
	gctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	gctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	created, err := repo.Create(gctx, cal)
	require.NoError(t, err)
	require.NotNil(t, created.ID)

	// GetByID
	got, err := repo.GetByID(gctx, *created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ProdID, got.ProdID)

	// Update
	created.ProdID = "-//Atomic Blend Updated//EN"
	updated, err := repo.Update(gctx, *created.ID, created)
	require.NoError(t, err)
	assert.Equal(t, "-//Atomic Blend Updated//EN", updated.ProdID)

	// GetAll
	list, total, err := repo.GetAll(gctx, userID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)

	// GetSince (use a time before creation)
	since := time.Now().Add(-time.Hour)
	sinceList, totalSince, err := repo.GetSince(gctx, userID, since, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, totalSince, int64(1))
	assert.NotEmpty(t, sinceList)

	// Delete
	err = repo.Delete(gctx, *created.ID)
	require.NoError(t, err)

	// Ensure deletion
	gotAfter, err := repo.GetByID(gctx, *created.ID)
	require.NoError(t, err)
	assert.Nil(t, gotAfter)
}
