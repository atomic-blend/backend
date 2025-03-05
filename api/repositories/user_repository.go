package repositories

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"atomic_blend_api/models"
	"atomic_blend_api/utils/db"
)

const userCollection = "users"

// UserRepository provides methods to interact with user data in the database
type UserRepository struct {
	collection *mongo.Collection
}


// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(database *mongo.Database) *UserRepository {
	if database == nil {
		database = db.Database
	}
	return &UserRepository{
		collection: database.Collection(userCollection),
	}
}

// Create adds a new user to the database
func (r *UserRepository) Create(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error) {
	// Generate an ID if not provided
	if user.ID == nil {
		id := primitive.NewObjectID()
		user.ID = &id
	}

	// Insert into database
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.UserEntity, error) {
	var user models.UserEntity
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update modifies an existing user in the database
func (r *UserRepository) Update(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error) {
	if user.ID == nil {
		return nil, errors.New("user ID is required for update")
	}

	filter := bson.M{"id": user.ID}
	update := bson.M{"$set": user}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Delete removes a user from the database by ID
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// FindByEmail finds a user by their email address
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.UserEntity, error) {
	var user models.UserEntity
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx *gin.Context, d primitive.ObjectID) (*models.UserEntity, error) {
	var user models.UserEntity
	err := r.collection.FindOne(ctx, bson.M{"_id": d}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}