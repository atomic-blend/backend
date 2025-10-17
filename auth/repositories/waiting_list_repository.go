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
	Create(ctx context.Context, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error)
	GetAll(ctx context.Context) ([]*waitinglist.WaitingList, error)
	GetByID(ctx context.Context, id string) (*waitinglist.WaitingList, error)
	Update(ctx context.Context, id string, waitingList *waitinglist.WaitingList) (*waitinglist.WaitingList, error)
	Delete(ctx context.Context, id string) error
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
