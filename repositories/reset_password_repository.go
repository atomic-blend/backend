package repositories

import (
	"atomic_blend_api/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserResetPasswordRequestRepositoryInterface defines the methods for the UserResetPasswordRequestRepository
type UserResetPasswordRequestRepositoryInterface interface {
	Create(ctx context.Context, request *models.UserResetPassword) (*models.UserResetPassword, error)
	FindByResetCode(ctx context.Context, resetCode string) (*models.UserResetPassword, error)
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) (*models.UserResetPassword, error)
}

// UserResetPasswordRequestRepository is the implementation of UserResetPasswordRequestRepositoryInterface
type UserResetPasswordRequestRepository struct {
	collection *mongo.Collection
}

// Ensure UserResetPasswordRequestRepository implements UserResetPasswordRequestRepositoryInterface
var _ UserResetPasswordRequestRepositoryInterface = (*UserResetPasswordRequestRepository)(nil)

// NewUserResetPasswordRequestRepository creates a new UserResetPasswordRequestRepository instance
func NewUserResetPasswordRequestRepository(db *mongo.Database) *UserResetPasswordRequestRepository {
	return &UserResetPasswordRequestRepository{
		collection: db.Collection("user_reset_password_requests"),
	}
}

// Create inserts a new reset password request into the database
func (r *UserResetPasswordRequestRepository) Create(ctx context.Context, request *models.UserResetPassword) (*models.UserResetPassword, error) {
	// Set the created and updated timestamps
	request.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	request.UpdatedAt = request.CreatedAt

	// Insert the request into the database
	_, err := r.collection.InsertOne(ctx, request)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// FindByResetCode retrieves a reset password request by its reset code
func (r *UserResetPasswordRequestRepository) FindByResetCode(ctx context.Context, resetCode string) (*models.UserResetPassword, error) {
	filter := bson.M{"reset_code": resetCode}
	var request models.UserResetPassword
	err := r.collection.FindOne(ctx, filter).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other error
	}

	return &request, nil
}

// Delete removes a reset password request by its ID
func (r *UserResetPasswordRequestRepository) Delete(ctx context.Context, id string) error {
	// Convert the string ID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Delete the request from the database
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}

// FindByUserID retrieves a reset password request by its user ID
func (r *UserResetPasswordRequestRepository) FindByUserID(ctx context.Context, userID string) (*models.UserResetPassword, error) {
	// Convert the string user ID to an ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var request models.UserResetPassword
	err = r.collection.FindOne(ctx, filter).Decode(&request)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No document found
		}
		return nil, err // Other error
	}

	return &request, nil
}
