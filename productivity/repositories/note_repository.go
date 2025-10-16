package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/atomic-blend/backend/productivity/models"
	patchmodels "github.com/atomic-blend/backend/productivity/models/patch_models"
	keyconverter "github.com/atomic-blend/backend/shared/utils/key_converter"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const noteCollection = "notes"

// NoteRepositoryInterface defines the interface for note repository operations
type NoteRepositoryInterface interface {
	// GetAll retrieves notes for a user. If page and limit are both provided and >0, returns paginated results and total count. If either is nil or <=0, returns all notes and total count.
	GetAll(ctx context.Context, userID *primitive.ObjectID, page, limit *int64) ([]*models.NoteEntity, int64, error)
	GetByID(ctx context.Context, id string) (*models.NoteEntity, error)
	Create(ctx context.Context, note *models.NoteEntity) (*models.NoteEntity, error)
	Update(ctx context.Context, id string, note *models.NoteEntity) (*models.NoteEntity, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error
	UpdatePatch(ctx context.Context, patch *patchmodels.Patch) (*models.NoteEntity, error)
}

// NoteRepository handles database operations related to notes
type NoteRepository struct {
	collection *mongo.Collection
}

// NewNoteRepository creates a new note repository instance
func NewNoteRepository(db *mongo.Database) NoteRepositoryInterface {
	return &NoteRepository{
		collection: db.Collection("notes"),
	}
}

// GetAll retrieves notes for a user. If page and limit are both provided and >0, returns paginated results and total count. If either is nil or <=0, returns all notes and total count.
func (r *NoteRepository) GetAll(ctx context.Context, userID *primitive.ObjectID, page, limit *int64) ([]*models.NoteEntity, int64, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user"] = userID
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build find options: always sort by created_at desc to return most recent first
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	// Only apply pagination if both page and limit are provided and > 0
	if page != nil && limit != nil && *page > 0 && *limit > 0 {
		skip := (*page - 1) * *limit
		findOpts.SetSkip(skip)
		findOpts.SetLimit(*limit)
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notes []*models.NoteEntity
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, 0, err
	}

	return notes, totalCount, nil
}

// GetByID retrieves a note by its ID
func (r *NoteRepository) GetByID(ctx context.Context, id string) (*models.NoteEntity, error) {
	if id == "" {
		return nil, errors.New("note ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	var note models.NoteEntity

	err = r.collection.FindOne(ctx, filter).Decode(&note)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &note, nil
}

// Create creates a new note
func (r *NoteRepository) Create(ctx context.Context, note *models.NoteEntity) (*models.NoteEntity, error) {
	if note == nil {
		return nil, errors.New("note cannot be nil")
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	note.CreatedAt = now
	note.UpdatedAt = now

	// Generate a new ObjectID if not provided
	if note.ID == nil {
		id := primitive.NewObjectID()
		note.ID = &id
	}

	_, err := r.collection.InsertOne(ctx, note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

// Update updates an existing note
func (r *NoteRepository) Update(ctx context.Context, id string, note *models.NoteEntity) (*models.NoteEntity, error) {
	if id == "" {
		return nil, errors.New("note ID is required")
	}
	if note == nil {
		return nil, errors.New("note cannot be nil")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Set updated timestamp
	updatedAt := primitive.NewDateTimeFromTime(time.Now())

	// Prepare update document excluding the ID and created_at
	update := bson.M{"$set": bson.M{
		"title":      note.Title,
		"content":    note.Content,
		"user":       note.User,
		"deleted":    note.Deleted,
		"updated_at": updatedAt,
	}}

	filter := bson.M{"_id": objID}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// Return the updated note
	return r.GetByID(ctx, id)
}

// Delete deletes a note by its ID
func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("note ID is required")
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("note not found")
	}

	return nil
}

// DeleteByUserID deletes all notes for a specific user
func (r *NoteRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"user": userID}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// UpdatePatch updates a note based on a patch
func (r *NoteRepository) UpdatePatch(ctx context.Context, patch *patchmodels.Patch) (*models.NoteEntity, error) {
	if patch.Action != "update" {
		return nil, errors.New("only update action is supported")
	}

	if patch.ItemType != patchmodels.ItemTypeNote {
		return nil, errors.New("item type not supported")
	}

	updatePayload := bson.M{}
	for _, change := range patch.Changes {
		//convert Key from camelCase to snake_case
		key := keyconverter.ToSnakeCase(change.Key)
		value := change.Value
		if isNoteDateTime(key) {
			if dateValue, err := convertToNoteDateTime(change.Value, isNoteDateTimePointer(key)); err == nil {
				value = dateValue
			} else {
				return nil, errors.New("invalid date format for field: " + key)
			}
		} else if isNoteBooleanField(key) {
			if boolValue, err := convertToNoteBoolean(change.Value, isNoteBooleanPointer(key)); err == nil {
				value = boolValue
			} else {
				return nil, errors.New("invalid boolean format for field: " + key)
			}
		}

		updatePayload[key] = value
	}

	updatePayload["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": patch.ItemID}, bson.M{"$set": updatePayload})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, patch.ItemID.Hex())
}

// Helper function to check if a note field is a date/time field
func isNoteDateTime(fieldName string) bool {
	dateTimeFields := []string{"created_at", "updated_at"}
	for _, field := range dateTimeFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to check if a field should be a pointer to DateTime
func isNoteDateTimePointer(fieldName string) bool {
	return false // Note entity doesn't have any dateTime pointer fields
}

// Helper function to check if a field is a boolean field
func isNoteBooleanField(fieldName string) bool {
	booleanFields := []string{"deleted"}
	for _, field := range booleanFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to check if a field should be a pointer to bool
func isNoteBooleanPointer(fieldName string) bool {
	pointerFields := []string{"deleted"}
	for _, field := range pointerFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to convert various date formats to primitive.DateTime
func convertToNoteDateTime(value interface{}, isPointer bool) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	var parsedTime time.Time
	var err error

	switch v := value.(type) {
	case string:
		// Try parsing different date formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05",
		}

		for _, format := range formats {
			if parsedTime, err = time.Parse(format, v); err == nil {
				break
			}
		}
		if err != nil {
			return nil, errors.New("unable to parse date string")
		}
	case int64:
		// Unix timestamp in milliseconds
		parsedTime = time.Unix(0, v*int64(time.Millisecond))
	case float64:
		// Unix timestamp in seconds (JavaScript timestamps are often float64)
		parsedTime = time.Unix(int64(v), 0)
	default:
		return nil, errors.New("unsupported date format")
	}

	dt := primitive.NewDateTimeFromTime(parsedTime)

	if isPointer {
		return &dt, nil
	}
	return dt, nil
}

// Helper function to convert various boolean formats to bool
func convertToNoteBoolean(value interface{}, isPointer bool) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	var boolValue bool

	switch v := value.(type) {
	case bool:
		boolValue = v
	case string:
		switch v {
		case "true", "True", "TRUE", "1":
			boolValue = true
		case "false", "False", "FALSE", "0":
			boolValue = false
		default:
			return nil, errors.New("unable to parse boolean string")
		}
	case int:
		boolValue = v != 0
	case int64:
		boolValue = v != 0
	case float64:
		boolValue = v != 0
	default:
		return nil, errors.New("unsupported boolean format")
	}

	if isPointer {
		return &boolValue, nil
	}
	return boolValue, nil
}
