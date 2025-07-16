package repositories

import (
	"atomic-blend/backend/productivity/models"
	"atomic-blend/backend/productivity/utils/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const timeEntryCollection = "time_entries"

// TimeEntryRepositoryInterface defines methods that a TimeEntryRepository must implement
type TimeEntryRepositoryInterface interface {
	Create(ctx context.Context, timeEntry *models.TimeEntry) (*models.TimeEntry, error)
	GetByID(ctx context.Context, id string) (*models.TimeEntry, error)
	GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TimeEntry, error)
	Update(ctx context.Context, id string, timeEntry *models.TimeEntry) (*models.TimeEntry, error)
	Delete(ctx context.Context, id string) error
}

// TimeEntryRepository provides methods to interact with time entry data in the database
type TimeEntryRepository struct {
	collection *mongo.Collection
}

// Ensure TimeEntryRepository implements TimeEntryRepositoryInterface
var _ TimeEntryRepositoryInterface = (*TimeEntryRepository)(nil)

// NewTimeEntryRepository creates a new instance of TimeEntryRepository
func NewTimeEntryRepository(database *mongo.Database) *TimeEntryRepository {
	if database == nil {
		database = db.Database
	}
	return &TimeEntryRepository{
		collection: database.Collection(timeEntryCollection),
	}
}

// Create adds a new time entry to the database
func (r *TimeEntryRepository) Create(ctx context.Context, timeEntry *models.TimeEntry) (*models.TimeEntry, error) {
	// Generate new ObjectID if not provided
	if timeEntry.ID == nil {
		id := primitive.NewObjectID()
		timeEntry.ID = &id
	}

	// Set creation timestamp
	now := time.Now().Format(time.RFC3339)
	timeEntry.CreatedAt = now
	timeEntry.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, timeEntry)
	if err != nil {
		return nil, err
	}

	return timeEntry, nil
}

// GetByID retrieves a time entry by its ID
func (r *TimeEntryRepository) GetByID(ctx context.Context, id string) (*models.TimeEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var timeEntry models.TimeEntry
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&timeEntry)
	if err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// GetAll retrieves all time entries with optional user filtering
func (r *TimeEntryRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TimeEntry, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var timeEntries []*models.TimeEntry
	if err := cursor.All(ctx, &timeEntries); err != nil {
		return nil, err
	}

	return timeEntries, nil
}

// Update modifies an existing time entry in the database
func (r *TimeEntryRepository) Update(ctx context.Context, id string, timeEntry *models.TimeEntry) (*models.TimeEntry, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Set update timestamp
	timeEntry.UpdatedAt = time.Now().Format(time.RFC3339)

	// Create update document excluding the ID field
	updateDoc := bson.M{
		"$set": bson.M{
			"start_date": timeEntry.StartDate,
			"end_date":   timeEntry.EndDate,
			"duration":   timeEntry.Duration,
			"timer":      timeEntry.Timer,
			"pomodoro":   timeEntry.Pomodoro,
			"updated_at": timeEntry.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, updateDoc)
	if err != nil {
		return nil, err
	}

	// Return the updated time entry
	return r.GetByID(ctx, id)
}

// Delete removes a time entry from the database
func (r *TimeEntryRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
