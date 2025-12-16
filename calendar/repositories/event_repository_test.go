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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupEventTest(t *testing.T) (EventRepositoryInterface, func()) {
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	db := client.Database("test_db")
	repo := NewEventRepository(db)

	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestEvent() *models.Event {
	now := primitive.NewDateTimeFromTime(time.Now())
	uid := "evt-12345"
	summary := "Test Event"
	return &models.Event{
		UID:       uid,
		DTStamp:   time.Now(),
		Summary:   summary,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func TestEventRepository_CreateGetUpdateDelete(t *testing.T) {
	repo, cleanup := setupEventTest(t)
	defer cleanup()

	calendarID := primitive.NewObjectID()

	ev := createTestEvent()
	ev.CalendarID = &calendarID

	gctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	gctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	created, err := repo.Create(gctx, ev)
	require.NoError(t, err)
	require.NotNil(t, created)

	// Find inserted _id via the concrete repository collection
	er, ok := repo.(*EventRepository)
	require.True(t, ok)

	cctx := gctx.Request.Context()
	var doc bson.M
	err = er.collection.FindOne(cctx, bson.M{"uid": created.UID}).Decode(&doc)
	require.NoError(t, err)
	id, ok := doc["_id"].(primitive.ObjectID)
	require.True(t, ok)

	// GetByID
	got, err := repo.GetByID(gctx, id)
	require.NoError(t, err)
	assert.Equal(t, created.Summary, got.Summary)

	// Update
	newSummary := "Updated Event"
	created.Summary = newSummary
	updated, err := repo.Update(gctx, id, created)
	require.NoError(t, err)
	assert.Equal(t, newSummary, updated.Summary)

	// GetAll
	list, total, err := repo.GetAll(gctx, calendarID, 1, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)

	// GetSince
	since := time.Now().Add(-time.Hour)
	sinceList, totalSince, err := repo.GetSince(gctx, calendarID, since, 1, 10)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, totalSince, int64(1))
	assert.NotEmpty(t, sinceList)

	// GetByUID
	gotByUID, err := repo.GetByUID(gctx, created.UID)
	require.NoError(t, err)
	assert.Equal(t, created.UID, gotByUID.UID)

	// Delete
	err = repo.Delete(gctx, id)
	require.NoError(t, err)

	// Ensure deletion
	gotAfter, err := repo.GetByID(gctx, id)
	require.NoError(t, err)
	assert.Nil(t, gotAfter)
}
