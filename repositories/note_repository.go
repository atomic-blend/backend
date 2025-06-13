package repositories

import (
	"atomic_blend_api/models"
	"context"
	"errors"
	"time"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const noteCollection = "notes"

// NoteRepositoryInterface defines the interface for note repository operations
type NoteRepositoryInterface interface {
	GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.NoteEntity, error)
	GetByID(ctx context.Context, id string) (*models.NoteEntity, error)
	Create(ctx context.Context, note *models.NoteEntity) (*models.NoteEntity, error)
	Update(ctx context.Context, id string, note *models.NoteEntity) (*models.NoteEntity, error)
	Delete(ctx context.Context, id string) error
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

// GetAll retrieves all notes with optional user filtering
func (r *NoteRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.NoteEntity, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []*models.NoteEntity
	for cursor.Next(ctx) {
		var note models.NoteEntity
		if err := cursor.Decode(&note); err != nil {
			return nil, err
		}
		notes = append(notes, &note)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notes, nil
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
