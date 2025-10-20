// Package repositories provides data access layer implementations for the auth service.
package repositories

import (
	"context"
	"time"

	waitinglist "github.com/atomic-blend/backend/auth/models/waiting_list"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// waitingListCollection is the name of the collection in the database
const waitingListCollection = "waiting_list"

// WaitingListRepositoryInterface defines the interface for waiting list repository operations
type WaitingListRepositoryInterface interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error)
	GetAll(ctx context.Context) ([]*waitinglist.WaitingList, error)
	GetByID(ctx context.Context, id string) (*waitinglist.WaitingList, error)
	GetByEmail(ctx context.Context, email string) (*waitinglist.WaitingList, error)
	GetByCode(ctx context.Context, code string) (*waitinglist.WaitingList, error)
	GetPositionByEmail(ctx context.Context, email string) (int64, error)
	Update(ctx context.Context, id string, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error)
	Delete(ctx context.Context, id string) error
	DeleteByEmail(ctx context.Context, email string) error
	DeleteByCode(ctx context.Context, code string) error
}

// WaitingListRepository handles database operations related to waiting lists
type WaitingListRepository struct {
	collection *mongo.Collection
}

// NewWaitingListRepository creates a new waiting list repository instance
func NewWaitingListRepository(database *mongo.Database) WaitingListRepositoryInterface {
	return &WaitingListRepository{
		collection: database.Collection(waitingListCollection),
	}
}

// Count returns the number of waiting list records
func (r *WaitingListRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Create creates a new waiting list
func (r *WaitingListRepository) Create(ctx context.Context, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	if waitingList.ID == nil {
		id := primitive.NewObjectID()
		waitingList.ID = &id
	}

	waitingList.CreatedAt = &now
	waitingList.UpdatedAt = &now

	_, err := r.collection.InsertOne(ctx, waitingList)
	if err != nil {
		return nil, err
	}

	return waitingList, nil
}

// GetAll retrieves all waiting lists sorted from oldest to newest by created_at
func (r *WaitingListRepository) GetAll(ctx context.Context) ([]*waitinglist.WaitingList, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var waitingLists []*waitinglist.WaitingList
	if err := cursor.All(ctx, &waitingLists); err != nil {
		return nil, err
	}

	return waitingLists, nil
}

// GetByID retrieves a waiting list by its ID
func (r *WaitingListRepository) GetByID(ctx context.Context, id string) (*waitinglist.WaitingList, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var waitingList waitinglist.WaitingList
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&waitingList)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &waitingList, nil
}

// GetByEmail retrieves a waiting list by email address
func (r *WaitingListRepository) GetByEmail(ctx context.Context, email string) (*waitinglist.WaitingList, error) {
	var waitingList waitinglist.WaitingList
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&waitingList)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &waitingList, nil
}

// GetByCode retrieves a waiting list by code
func (r *WaitingListRepository) GetByCode(ctx context.Context, code string) (*waitinglist.WaitingList, error) {
	var waitingList waitinglist.WaitingList
	err := r.collection.FindOne(ctx, bson.M{"code": code}).Decode(&waitingList)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &waitingList, nil
}

// Update modifies an existing waiting list in the database
func (r *WaitingListRepository) Update(ctx context.Context, id string, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	waitingList.UpdatedAt = &now

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": waitingList})
	if err != nil {
		return nil, err
	}

	return waitingList, nil
}

// Delete removes a waiting list from the database by its ID
func (r *WaitingListRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	return nil
}

// DeleteByEmail removes a waiting list from the database by email address
func (r *WaitingListRepository) DeleteByEmail(ctx context.Context, email string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		return err
	}

	return nil
}

// DeleteByCode removes a waiting list from the database by code
func (r *WaitingListRepository) DeleteByCode(ctx context.Context, code string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"code": code})
	if err != nil {
		return err
	}

	return nil
}

// GetPositionByEmail returns the position of a waiting list record by email (0-based index)
// Position is calculated by counting records created before this record (sorted by created_at)
func (r *WaitingListRepository) GetPositionByEmail(ctx context.Context, email string) (int64, error) {
	// First get the record to find its created_at timestamp
	var waitingList waitinglist.WaitingList
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&waitingList)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil // Return 0 if not found
		}
		return 0, err
	}

	// Count records created before this record (created_at < current record's created_at)
	count, err := r.collection.CountDocuments(ctx, bson.M{
		"created_at": bson.M{"$lt": waitingList.CreatedAt},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}
