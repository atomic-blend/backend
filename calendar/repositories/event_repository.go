package repositories

import (
    "context"
    "time"

    "github.com/atomic-blend/backend/calendar/models"
    "github.com/atomic-blend/backend/shared/utils/db"
    "github.com/gin-gonic/gin"

    bson "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

const eventCollection = "events"

// EventRepositoryInterface defines event repository operations
type EventRepositoryInterface interface {
    GetAll(ctx *gin.Context, calendarID primitive.ObjectID, page int64, size int64) ([]*models.Event, int64, error)
    GetSince(ctx *gin.Context, calendarID primitive.ObjectID, since time.Time, page int64, size int64) ([]*models.Event, int64, error)
    GetByID(ctx *gin.Context, id primitive.ObjectID) (*models.Event, error)
    GetByUID(ctx *gin.Context, uid string) (*models.Event, error)
    GetByUIDWithContext(ctx context.Context, uid string) (*models.Event, error)
    Create(ctx *gin.Context, event *models.Event) (*models.Event, error)
    CreateWithContext(ctx context.Context, event *models.Event) (*models.Event, error)
    Update(ctx *gin.Context, id primitive.ObjectID, event *models.Event) (*models.Event, error)
    Delete(ctx *gin.Context, id primitive.ObjectID) error
}

// EventRepository handles DB ops for events
type EventRepository struct {
    collection *mongo.Collection
}

// NewEventRepository creates a new repository instance
func NewEventRepository(database *mongo.Database) EventRepositoryInterface {
    if database == nil {
        database = db.Database
    }
    return &EventRepository{collection: database.Collection(eventCollection)}
}

// GetAll retrieves events for a calendar with pagination
func (r *EventRepository) GetAll(ctx *gin.Context, calendarID primitive.ObjectID, page int64, size int64) ([]*models.Event, int64, error) {
    filter := bson.M{"calendar_id": calendarID}

    findOptions := options.Find()
    if size > 0 {
        findOptions.SetLimit(size)
        if page > 1 {
            findOptions.SetSkip((page - 1) * size)
        }
    }

    cctx := ctx.Request.Context()

    cursor, err := r.collection.Find(cctx, filter, findOptions)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(cctx)

    var events []*models.Event
    if err := cursor.All(cctx, &events); err != nil {
        return nil, 0, err
    }

    totalCount, _ := r.collection.CountDocuments(cctx, filter)

    return events, totalCount, nil
}

// GetSince retrieves events updated since a time with pagination
func (r *EventRepository) GetSince(ctx *gin.Context, calendarID primitive.ObjectID, since time.Time, page int64, size int64) ([]*models.Event, int64, error) {
    filter := bson.M{"calendar_id": calendarID, "updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(since)}}

    findOptions := options.Find()
    if size > 0 {
        findOptions.SetLimit(size)
        if page > 1 {
            findOptions.SetSkip((page - 1) * size)
        }
    }

    cctx := ctx.Request.Context()

    cursor, err := r.collection.Find(cctx, filter, findOptions)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(cctx)

    var events []*models.Event
    if err := cursor.All(cctx, &events); err != nil {
        return nil, 0, err
    }

    totalCount, _ := r.collection.CountDocuments(cctx, filter)

    return events, totalCount, nil
}

// GetByID returns an event by its Mongo _id
func (r *EventRepository) GetByID(ctx *gin.Context, id primitive.ObjectID) (*models.Event, error) {
    filter := bson.M{"_id": id}
    var event models.Event
    cctx := ctx.Request.Context()
    err := r.collection.FindOne(cctx, filter).Decode(&event)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &event, nil
}

// Create inserts a new event
func (r *EventRepository) Create(ctx *gin.Context, event *models.Event) (*models.Event, error) {
    now := primitive.NewDateTimeFromTime(time.Now())
    event.CreatedAt = &now
    event.UpdatedAt = &now

    cctx := ctx.Request.Context()
    _, err := r.collection.InsertOne(cctx, event)
    if err != nil {
        return nil, err
    }
    return event, nil
}

// Update modifies an existing event by _id
func (r *EventRepository) Update(ctx *gin.Context, id primitive.ObjectID, event *models.Event) (*models.Event, error) {
    now := primitive.NewDateTimeFromTime(time.Now())
    event.UpdatedAt = &now

    filter := bson.M{"_id": id}
    update := bson.M{"$set": event}

    cctx := ctx.Request.Context()
    _, err := r.collection.UpdateOne(cctx, filter, update)
    if err != nil {
        return nil, err
    }

    return r.GetByID(ctx, id)
}

// Delete removes an event document
func (r *EventRepository) Delete(ctx *gin.Context, id primitive.ObjectID) error {
    filter := bson.M{"_id": id}
    cctx := ctx.Request.Context()
    _, err := r.collection.DeleteOne(cctx, filter)
    return err
}

// CreateWithContext inserts a new event using the provided context
func (r *EventRepository) CreateWithContext(ctx context.Context, event *models.Event) (*models.Event, error) {
    now := primitive.NewDateTimeFromTime(time.Now())
    event.CreatedAt = &now
    event.UpdatedAt = &now

    _, err := r.collection.InsertOne(ctx, event)
    if err != nil {
        return nil, err
    }
    return event, nil
}

// GetByUID retrieves an event by its UID
func (r *EventRepository) GetByUID(ctx *gin.Context, uid string) (*models.Event, error) {
    cctx := ctx.Request.Context()
    return r.GetByUIDWithContext(cctx, uid)
}

// GetByUIDWithContext retrieves an event by its UID using provided context
func (r *EventRepository) GetByUIDWithContext(ctx context.Context, uid string) (*models.Event, error) {
    filter := bson.M{"uid": uid}
    var event models.Event
    err := r.collection.FindOne(ctx, filter).Decode(&event)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }
    return &event, nil
}
