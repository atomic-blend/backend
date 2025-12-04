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

const calendarCollection = "calendars"

// CalendarRepositoryInterface defines calendar repository operations
type CalendarRepositoryInterface interface {
	GetAll(ctx *gin.Context, userID primitive.ObjectID, page int64, size int64) ([]*models.Calendar, int64, error)
	GetSince(ctx *gin.Context, userID primitive.ObjectID, since time.Time, page int64, size int64) ([]*models.Calendar, int64, error)
	GetByID(ctx *gin.Context, id primitive.ObjectID) (*models.Calendar, error)
	GetByName(ctx *gin.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error)
	GetByNameWithContext(ctx context.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error)
	Create(ctx *gin.Context, calendar *models.Calendar) (*models.Calendar, error)
	CreateWithContext(ctx context.Context, calendar *models.Calendar) (*models.Calendar, error)
	Update(ctx *gin.Context, id primitive.ObjectID, calendar *models.Calendar) (*models.Calendar, error)
	Delete(ctx *gin.Context, id primitive.ObjectID) error
}

// CalendarRepository handles DB ops for calendars
type CalendarRepository struct {
	collection *mongo.Collection
}

// NewCalendarRepository creates a new repository instance
func NewCalendarRepository(database *mongo.Database) CalendarRepositoryInterface {
	if database == nil {
		database = db.Database
	}
	return &CalendarRepository{collection: database.Collection(calendarCollection)}
}

// GetAll retrieves calendars for a user with pagination
func (r *CalendarRepository) GetAll(ctx *gin.Context, userID primitive.ObjectID, page int64, size int64) ([]*models.Calendar, int64, error) {
	filter := bson.M{"user_id": userID}

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

	var calendars []*models.Calendar
	if err := cursor.All(cctx, &calendars); err != nil {
		return nil, 0, err
	}

	totalCount, _ := r.collection.CountDocuments(cctx, filter)

	return calendars, totalCount, nil
}

// GetSince retrieves calendars updated since a time with pagination
func (r *CalendarRepository) GetSince(ctx *gin.Context, userID primitive.ObjectID, since time.Time, page int64, size int64) ([]*models.Calendar, int64, error) {
	filter := bson.M{"user_id": userID, "updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(since)}}

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

	var calendars []*models.Calendar
	if err := cursor.All(cctx, &calendars); err != nil {
		return nil, 0, err
	}

	totalCount, _ := r.collection.CountDocuments(cctx, filter)

	return calendars, totalCount, nil
}

// GetByID returns a calendar by its ID
func (r *CalendarRepository) GetByID(ctx *gin.Context, id primitive.ObjectID) (*models.Calendar, error) {
	filter := bson.M{"_id": id}
	var calendar models.Calendar
	cctx := ctx.Request.Context()
	err := r.collection.FindOne(cctx, filter).Decode(&calendar)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &calendar, nil
}

// Create inserts a new calendar
func (r *CalendarRepository) Create(ctx *gin.Context, calendar *models.Calendar) (*models.Calendar, error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	if calendar.ID == nil {
		id := primitive.NewObjectID()
		calendar.ID = &id
	}
	calendar.CreatedAt = &now
	calendar.UpdatedAt = &now

	cctx := ctx.Request.Context()
	_, err := r.collection.InsertOne(cctx, calendar)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}

// Update modifies an existing calendar
func (r *CalendarRepository) Update(ctx *gin.Context, id primitive.ObjectID, calendar *models.Calendar) (*models.Calendar, error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	calendar.UpdatedAt = &now

	filter := bson.M{"_id": id}
	update := bson.M{"$set": calendar}

	cctx := ctx.Request.Context()
	_, err := r.collection.UpdateOne(cctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

// Delete removes a calendar document
func (r *CalendarRepository) Delete(ctx *gin.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	cctx := ctx.Request.Context()
	_, err := r.collection.DeleteOne(cctx, filter)
	return err
}

// CreateWithContext inserts a new calendar using the provided context
func (r *CalendarRepository) CreateWithContext(ctx context.Context, calendar *models.Calendar) (*models.Calendar, error) {
	now := primitive.NewDateTimeFromTime(time.Now())
	if calendar.ID == nil {
		id := primitive.NewObjectID()
		calendar.ID = &id
	}
	calendar.CreatedAt = &now
	calendar.UpdatedAt = &now

	_, err := r.collection.InsertOne(ctx, calendar)
	if err != nil {
		return nil, err
	}
	return calendar, nil
}

// GetByName retrieves a calendar by user ID and name
func (r *CalendarRepository) GetByName(ctx *gin.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error) {
	cctx := ctx.Request.Context()
	return r.GetByNameWithContext(cctx, userID, name)
}

// GetByNameWithContext retrieves a calendar by user ID and name using the provided context
func (r *CalendarRepository) GetByNameWithContext(ctx context.Context, userID primitive.ObjectID, name *string) (*models.Calendar, error) {
	filter := bson.M{"user_id": userID, "name": name}
	var calendar models.Calendar
	err := r.collection.FindOne(ctx, filter).Decode(&calendar)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &calendar, nil
}
