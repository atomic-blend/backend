package repositories

import (
	"atomic_blend_api/models"
	"atomic_blend_api/utils/db"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const habitCollection = "habits"
const habitEntryCollection = "habit_entries"

// HabitRepositoryInterface defines methods that a HabitRepository must implement
type HabitRepositoryInterface interface {
	Create(ctx context.Context, habit *models.Habit) (*models.Habit, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Habit, error)
	GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) (*models.Habit, error)
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Habit Entry methods
	AddEntry(ctx context.Context, entry *models.HabitEntry) (*models.HabitEntry, error)
	GetEntriesByHabitID(ctx context.Context, habitID primitive.ObjectID) ([]models.HabitEntry, error)
	UpdateEntry(ctx context.Context, entry *models.HabitEntry) (*models.HabitEntry, error)
	DeleteEntry(ctx context.Context, id primitive.ObjectID) error
}

// HabitRepository provides methods to interact with habit data in the database
type HabitRepository struct {
	collection      *mongo.Collection
	entryCollection *mongo.Collection
}

// Ensure HabitRepository implements HabitRepositoryInterface
var _ HabitRepositoryInterface = (*HabitRepository)(nil)

// NewHabitRepository creates a new instance of HabitRepository
func NewHabitRepository(database *mongo.Database) *HabitRepository {
	if database == nil {
		database = db.Database
	}
	return &HabitRepository{
		collection:      database.Collection(habitCollection),
		entryCollection: database.Collection(habitEntryCollection),
	}
}

// Create adds a new habit to the database
func (r *HabitRepository) Create(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	// Generate an ID if not provided
	if habit.ID.IsZero() {
		habit.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now().Format(time.RFC3339)
	habit.CreatedAt = &now
	habit.UpdatedAt = &now

	// Insert into database
	_, err := r.collection.InsertOne(ctx, habit)
	if err != nil {
		return nil, err
	}
	return habit, nil
}

// GetByID retrieves a habit by its ID
func (r *HabitRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Habit, error) {
	var habit models.Habit
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&habit)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

// GetAll retrieves all habits for a specific user
func (r *HabitRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Habit, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user_id"] = *userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var habits []*models.Habit
	if err := cursor.All(ctx, &habits); err != nil {
		return nil, err
	}
	return habits, nil
}

// Update modifies an existing habit in the database
func (r *HabitRepository) Update(ctx context.Context, habit *models.Habit) (*models.Habit, error) {
	// Update timestamp
	now := time.Now().Format(time.RFC3339)
	habit.UpdatedAt = &now

	filter := bson.M{"_id": habit.ID}
	update := bson.M{"$set": habit}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return habit, nil
}

// Delete removes a habit from the database
func (r *HabitRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddEntry adds a new habit entry to the database
func (r *HabitRepository) AddEntry(ctx context.Context, entry *models.HabitEntry) (*models.HabitEntry, error) {
	// Generate an ID if not provided
	if entry.ID.IsZero() {
		entry.ID = primitive.NewObjectID()
	}

	// Set timestamps
	now := time.Now().Format(time.RFC3339)
	entry.CreatedAt = now
	entry.UpdatedAt = now

	// Insert into database
	_, err := r.entryCollection.InsertOne(ctx, entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// GetEntriesByHabitID retrieves all entries for a specific habit
func (r *HabitRepository) GetEntriesByHabitID(ctx context.Context, habitID primitive.ObjectID) ([]models.HabitEntry, error) {
	filter := bson.M{"habit_id": habitID}

	cursor, err := r.entryCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []models.HabitEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

// UpdateEntry modifies an existing habit entry in the database
func (r *HabitRepository) UpdateEntry(ctx context.Context, entry *models.HabitEntry) (*models.HabitEntry, error) {
	// Update timestamp
	now := time.Now().Format(time.RFC3339)
	entry.UpdatedAt = now

	filter := bson.M{"_id": entry.ID}
	update := bson.M{"$set": entry}

	_, err := r.entryCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

// DeleteEntry removes a habit entry from the database
func (r *HabitRepository) DeleteEntry(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.entryCollection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
