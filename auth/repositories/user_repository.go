package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/atomic-blend/backend/auth/models"
	"github.com/atomic-blend/backend/auth/utils/db"
	regexutils "github.com/atomic-blend/backend/auth/utils/regex"
)

const userCollection = "users"

// UserRepositoryInterface defines methods that a UserRepository must implement
type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error)
	GetByID(ctx context.Context, id string) (*models.UserEntity, error)
	Update(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error)
	Delete(ctx context.Context, id string) error
	FindByEmail(ctx context.Context, email string) (*models.UserEntity, error)
	FindByID(ctx *gin.Context, id primitive.ObjectID) (*models.UserEntity, error)
	ResetAllUserData(ctx *gin.Context, userID primitive.ObjectID) error
	AddPurchase(ctx *gin.Context, userID primitive.ObjectID, purchaseEntry *models.PurchaseEntity) error
}

// UserRepository provides methods to interact with user data in the database
type UserRepository struct {
	collection           *mongo.Collection
	taskCollection       *mongo.Collection
	habitCollection      *mongo.Collection
	habitEntryCollection *mongo.Collection
	tagCollection        *mongo.Collection
}

// Ensure UserRepository implements UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(database *mongo.Database) *UserRepository {
	if database == nil {
		database = db.Database
	}
	return &UserRepository{
		collection: database.Collection(userCollection),
		//TODO: replace with grpc calls
		// taskCollection:       database.Collection(taskCollection),
		// habitCollection:      database.Collection(habitCollection),
		// habitEntryCollection: database.Collection(habitEntryCollection),
		// tagCollection:        database.Collection(tagCollection),
	}
}

// GetAllIterable retrieves all users from the database
func (r *UserRepository) GetAllIterable(ctx context.Context) (*mongo.Cursor, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

// Create adds a new user to the database
func (r *UserRepository) Create(ctx context.Context, user *models.UserEntity) (*models.UserEntity, error) {
	// Generate an ID if not provided
	if user.ID == nil {
		id := primitive.NewObjectID()
		user.ID = &id
	}

	// Sanitize email if present
	if user.Email != nil {
		sanitizedEmail := regexutils.SanitizeString(*user.Email)
		user.Email = &sanitizedEmail
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	user.CreatedAt = &now
	user.UpdatedAt = &now
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
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

	// Sanitize email if present
	if user.Email != nil {
		sanitizedEmail := regexutils.SanitizeString(*user.Email)
		user.Email = &sanitizedEmail
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	user.UpdatedAt = &now

	filter := bson.M{"_id": user.ID}
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
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
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
	// Validate email format
	if !regexutils.IsValidEmail(email) {
		return nil, errors.New("invalid email format")
	}

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

// FindByID finds a user by their ObjectID
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

// ResetAllUserData deletes all personal data associated with a user
func (r *UserRepository) ResetAllUserData(ctx *gin.Context, userID primitive.ObjectID) error {
	// Delete all tasks for the user
	_, err := r.taskCollection.DeleteMany(ctx, bson.M{"user": userID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user tasks")
		return err
	}

	// Delete all habits for the user
	_, err = r.habitCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user habits")
		return err
	}

	// Delete all habit entries for the user
	_, err = r.habitEntryCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user habit entries")
		return err
	}

	// Delete all tags for the user
	_, err = r.tagCollection.DeleteMany(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete user tags")
		return err
	}

	return nil
}

// AddPurchase adds a purchase entry to the user
func (r *UserRepository) AddPurchase(ctx *gin.Context, userID primitive.ObjectID, purchaseEntry *models.PurchaseEntity) error {
	user, error := r.FindByID(ctx, userID)
	if error != nil {
		log.Error().Err(error).Msg("Failed to find user")
		return error
	}
	if user == nil {
		log.Error().Msg("User not found")
		return errors.New("user not found")
	}
	purchaseEntry.ID = primitive.NewObjectID()
	purchaseEntry.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	purchaseEntry.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Add the purchase entry to the user's purchases
	user.Purchases = append(user.Purchases, purchaseEntry)
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"purchases": user.Purchases}})
	if err != nil {
		log.Error().Err(err).Msg("Failed to update user with purchase entry")
		return err
	}

	return nil
}
