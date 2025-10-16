// Package repositories is a package that contains the repository for the microservice
package repositories

import (
	"context"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/utils/db"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const draftMailCollection = "draft_mails"

// DraftMailRepositoryInterface defines the interface for draft mail repository operations
type DraftMailRepositoryInterface interface {
	// GetAll retrieves draft mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
	GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.SendMail, int64, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.SendMail, error)
	Create(ctx context.Context, sendMail *models.SendMail) (*models.SendMail, error)
	Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.SendMail, error)
	Trash(ctx context.Context, id primitive.ObjectID) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	// GetSince retrieves draft mails where updated_at is after the specified time for a specific user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
	GetSince(ctx context.Context, userID primitive.ObjectID, since time.Time, page, limit int64) ([]*models.SendMail, int64, error)
}

// DraftMailRepository handles database operations related to draft mails
type DraftMailRepository struct {
	collection *mongo.Collection
}

// NewDraftMailRepository creates a new draft mail repository instance
func NewDraftMailRepository(database *mongo.Database) DraftMailRepositoryInterface {
	if database == nil {
		database = db.Database
	}
	return &DraftMailRepository{
		collection: database.Collection(draftMailCollection),
	}
}

// GetAll retrieves draft mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
func (r *DraftMailRepository) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.SendMail, int64, error) {
	filter := bson.M{"mail.user_id": userID}
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []*models.SendMail{}, 0, err
	}

	var opts []*options.FindOptions
	if page > 0 && limit > 0 {
		skip := (page - 1) * limit
		opts = append(opts, &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		})
	}

	// Sort by created_at descending
	sort := bson.D{{Key: "created_at", Value: -1}}
	if len(opts) > 0 {
		opts[0].Sort = sort
	} else {
		opts = append(opts, &options.FindOptions{Sort: sort})
	}

	var cursor *mongo.Cursor
	if len(opts) > 0 {
		cursor, err = r.collection.Find(ctx, filter, opts[0])
	} else {
		cursor, err = r.collection.Find(ctx, filter)
	}
	if err != nil {
		return []*models.SendMail{}, 0, err
	}
	defer cursor.Close(ctx)

	var sendMails []*models.SendMail
	if err = cursor.All(ctx, &sendMails); err != nil {
		return []*models.SendMail{}, 0, err
	}

	return sendMails, totalCount, nil
}

// GetByID retrieves a draft mail by its ID
func (r *DraftMailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.SendMail, error) {
	filter := bson.M{"_id": id}

	var sendMail models.SendMail
	err := r.collection.FindOne(ctx, filter).Decode(&sendMail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &sendMail, nil
}

// Create creates a new draft mail
func (r *DraftMailRepository) Create(ctx context.Context, sendMail *models.SendMail) (*models.SendMail, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	if sendMail.ID == primitive.NilObjectID {
		id := primitive.NewObjectID()
		sendMail.ID = id
	}

	sendMail.CreatedAt = &now
	sendMail.UpdatedAt = &now

	// Set default status if not provided
	if sendMail.SendStatus == "" {
		sendMail.SendStatus = models.SendStatusPending
	}

	// RetryCounter should be nil when email is created - will be handled by worker

	_, err := r.collection.InsertOne(ctx, sendMail)
	if err != nil {
		return nil, err
	}

	return sendMail, nil
}

// Update updates a draft mail by its ID
func (r *DraftMailRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) (*models.SendMail, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{"_id": id}

	// Merge the provided update with the updated_at timestamp
	updateDoc := bson.M{
		"$set": bson.M{
			"updated_at": now,
		},
	}

	// Add the provided update fields to the $set operation
	if setFields, exists := update["$set"]; exists {
		if setMap, ok := setFields.(bson.M); ok {
			for key, value := range setMap {
				updateDoc["$set"].(bson.M)[key] = value
			}
		}
	} else {
		// If no $set in the provided update, add all fields to $set
		for key, value := range update {
			updateDoc["$set"].(bson.M)[key] = value
		}
	}

	var sendMail models.SendMail
	err := r.collection.FindOneAndUpdate(ctx, filter, updateDoc, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&sendMail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &sendMail, nil
}

// Trash soft deletes a draft mail by marking it as trashed
func (r *DraftMailRepository) Trash(ctx context.Context, id primitive.ObjectID) error {
	now := primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"trashed":    true,
			"updated_at": now,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a draft mail by its ID
func (r *DraftMailRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// GetSince retrieves draft mails where updated_at is after the specified time for a specific user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all draft mails and total count.
func (r *DraftMailRepository) GetSince(ctx context.Context, userID primitive.ObjectID, since time.Time, page, limit int64) ([]*models.SendMail, int64, error) {
	filter := bson.M{
		"mail.user_id": userID,
		"updated_at":   bson.M{"$gt": primitive.NewDateTimeFromTime(since)},
	}

	// Count total documents matching the filter
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return []*models.SendMail{}, 0, err
	}

	// Build find options: always sort by updated_at desc to return most recent first
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "updated_at", Value: -1}})

	if page > 0 && limit > 0 {
		skip := (page - 1) * limit
		findOpts.SetSkip(skip)
		findOpts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return []*models.SendMail{}, 0, err
	}
	defer cursor.Close(ctx)

	var sendMails []*models.SendMail
	if err = cursor.All(ctx, &sendMails); err != nil {
		return []*models.SendMail{}, 0, err
	}

	return sendMails, totalCount, nil
}
