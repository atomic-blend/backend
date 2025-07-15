package repositories

import (
	"productivity/utils/db"
	"context"
	"productivity/models"
	"time"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const folderCollection = "folders"

// FolderRepositoryInterface defines the interface for folder repository operations
type FolderRepositoryInterface interface {
	GetAll(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error)
	Create(ctx context.Context, folder *models.Folder) (*models.Folder, error)
	Update(ctx context.Context, id primitive.ObjectID, folder *models.Folder) (*models.Folder, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

// FolderRepository handles database operations related to folders
type FolderRepository struct {
	collection *mongo.Collection
}

// NewFolderRepository creates a new folder repository instance
func NewFolderRepository(database *mongo.Database) FolderRepositoryInterface {
	if database == nil {
		database = db.Database
	}
	return &FolderRepository{
		collection: database.Collection(folderCollection),
	}
}

// GetAll retrieves all folders for a user
func (r *FolderRepository) GetAll(ctx context.Context, userID primitive.ObjectID) ([]*models.Folder, error) {
	filter := bson.M{"user_id": userID}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var folders []*models.Folder
	if err = cursor.All(ctx, &folders); err != nil {
		return nil, err
	}

	return folders, nil
}

// GetByID retrieves a folder by its ID
func (r *FolderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Folder, error) {
	filter := bson.M{"_id": id}

	var folder models.Folder
	err := r.collection.FindOne(ctx, filter).Decode(&folder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &folder, nil
}

// Create creates a new folder
func (r *FolderRepository) Create(ctx context.Context, folder *models.Folder) (*models.Folder, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	if folder.ID == nil {
		id := primitive.NewObjectID()
		folder.ID = &id
	}

	folder.CreatedAt = &now
	folder.UpdatedAt = &now

	_, err := r.collection.InsertOne(ctx, folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

// Delete removes a folder
func (r *FolderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// First, get all tasks that belong to this folder and remove the folder_id
	filter := bson.M{"folder_id": id}
	update := bson.M{"$set": bson.M{"folder_id": nil}}

	// Update tasks collection to remove folder reference
	_, err := r.collection.Database().Collection("tasks").UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	// Then delete the folder
	deleteFilter := bson.M{"_id": id}
	_, err = r.collection.DeleteOne(ctx, deleteFilter)
	return err
}

// Update modifies an existing folder
func (r *FolderRepository) Update(ctx context.Context, id primitive.ObjectID, folder *models.Folder) (*models.Folder, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": folder}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return folder, nil
}