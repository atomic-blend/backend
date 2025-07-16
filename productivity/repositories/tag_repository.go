package repositories

import (
	"atomic-blend/backend/productivity/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const tagCollection = "tags"

// TagRepositoryInterface defines the interface for tag repository operations
type TagRepositoryInterface interface {
	GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Tag, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error)
	Create(ctx context.Context, tag *models.Tag) (*models.Tag, error)
	Update(ctx context.Context, tag *models.Tag) (*models.Tag, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// TagRepository handles database operations related to tags
type TagRepository struct {
	collection *mongo.Collection
}

// Ensure TagRepository implements TagRepositoryInterface
var _ TagRepositoryInterface = (*TagRepository)(nil)

// NewTagRepository creates a new tag repository instance
func NewTagRepository(db *mongo.Database) *TagRepository {
	return &TagRepository{
		collection: db.Collection(tagCollection),
	}
}

// GetAll retrieves all tags with optional user filtering
func (r *TagRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.Tag, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user_id"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []*models.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// GetByID retrieves a tag by its ID
func (r *TagRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Tag, error) {
	var tag models.Tag
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&tag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &tag, nil
}

// Create adds a new tag to the database
func (r *TagRepository) Create(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	// Generate a new ID if not provided
	if tag.ID == nil {
		id := primitive.NewObjectID()
		tag.ID = &id
	}

	tag.CreatedAt = &now
	tag.UpdatedAt = &now

	_, err := r.collection.InsertOne(ctx, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// Update modifies an existing tag in the database
func (r *TagRepository) Update(ctx context.Context, tag *models.Tag) (*models.Tag, error) {
	if tag.ID == nil {
		return nil, errors.New("tag ID is required")
	}

	// Check if the tag exists first
	existingTag, err := r.GetByID(ctx, *tag.ID)
	if err != nil {
		return nil, err
	}
	if existingTag == nil {
		return nil, errors.New("tag not found")
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	tag.UpdatedAt = &now

	filter := bson.M{"_id": tag.ID}
	update := bson.M{"$set": tag}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, *tag.ID)
}

// Delete removes a tag from the database
func (r *TagRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
