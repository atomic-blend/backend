package repositories

import (
	"atomic_blend_api/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const taskCollection = "tasks"

// TaskRepositoryInterface defines the interface for task repository operations
type TaskRepositoryInterface interface {
	GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TaskEntity, error)
	GetByID(ctx context.Context, id string) (*models.TaskEntity, error)
	Create(ctx context.Context, task *models.TaskEntity) (*models.TaskEntity, error)
	Update(ctx context.Context, id string, task *models.TaskEntity) (*models.TaskEntity, error)
	Delete(ctx context.Context, id string) error
	AddTimeEntry(ctx context.Context, taskID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error)
	RemoveTimeEntry(ctx context.Context, taskID string, timeEntryID string) (*models.TaskEntity, error)
	UpdateTimeEntry(ctx context.Context, taskID string, timeEntryID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error)
}

// TaskRepository handles database operations related to tasks
type TaskRepository struct {
	collection *mongo.Collection
}

// NewTaskRepository creates a new task repository instance
func NewTaskRepository(db *mongo.Database) TaskRepositoryInterface {
	return &TaskRepository{
		collection: db.Collection("tasks"),
	}
}

// GetAll retrieves all tasks with optional user filtering
func (r *TaskRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.TaskEntity, error) {
	filter := bson.M{}
	if userID != nil {
		filter["user"] = userID
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.TaskEntity
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(ctx context.Context, id string) (*models.TaskEntity, error) {
	var task models.TaskEntity
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	err = r.collection.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &task, nil
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *models.TaskEntity) (*models.TaskEntity, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	// Generate a new ID if not provided
	if task.ID == "" {
		objID := primitive.NewObjectID()
		task.ID = objID.Hex()
	}

	task.CreatedAt = now
	task.UpdatedAt = now

	// Convert string ID to ObjectID for storing in MongoDB
	objID, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return nil, err
	}

	_, err = r.collection.InsertOne(ctx, bson.M{
		"_id":          objID,
		"title":        task.Title,
		"user":         task.User,
		"description":  task.Description,
		"start_date":   task.StartDate,
		"end_date":     task.EndDate,
		"completed":    task.Completed,
		"reminders":    task.Reminders,
		"priority":     task.Priority,
		"time_entries": task.TimeEntries,
		"created_at":   task.CreatedAt,
		"updated_at":   task.UpdatedAt,
	})

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepository) AddTimeEntry(ctx context.Context, taskID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	// set the array for time entries if not defined
	task := &models.TaskEntity{}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	if task.TimeEntries == nil {
		task.TimeEntries = []*models.TimeEntry{}
	}

	// Add the new time entry to the task
	task.TimeEntries = append(task.TimeEntries, timeEntry)
	update := bson.M{
		"$set": bson.M{
			"time_entries": task.TimeEntries,
		},
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, taskID)
}

// RemoveTimeEntry removes a time entry from a task
func (r *TaskRepository) RemoveTimeEntry(ctx context.Context, taskID string, timeEntryID string) (*models.TaskEntity, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	task := &models.TaskEntity{}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	if task.TimeEntries == nil {
		return nil, errors.New("no time entries found")
	}

	// Convert timeEntryID string to ObjectID
	timeEntryObjID, err := primitive.ObjectIDFromHex(timeEntryID)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$pull": bson.M{
			"time_entries": bson.M{"_id": timeEntryObjID},
		},
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, taskID)
}

// UpdateTimeEntry updates a time entry in a task
func (r *TaskRepository) UpdateTimeEntry(ctx context.Context, taskID string, timeEntryID string, timeEntry *models.TimeEntry) (*models.TaskEntity, error) {
	objID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		return nil, err
	}

	task := &models.TaskEntity{}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	if task.TimeEntries == nil {
		return nil, errors.New("no time entries found")
	}

	update := bson.M{
		"$set": bson.M{
			"time_entries.$[entry].start_date": timeEntry.StartDate,
			"time_entries.$[entry].end_date":   timeEntry.EndDate,
		},
	}

	filter := bson.M{"_id": objID}

	// Convert timeEntryID string to ObjectID
	timeEntryObjID, err := primitive.ObjectIDFromHex(timeEntryID)
	if err != nil {
		return nil, err
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"entry._id": timeEntryObjID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)
	_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, taskID)
}

// Update updates an existing task
func (r *TaskRepository) Update(ctx context.Context, id string, task *models.TaskEntity) (*models.TaskEntity, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	task.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	update := bson.M{
		"$set": bson.M{
			"title":       task.Title,
			"description": task.Description,
			"start_date":  task.StartDate,
			"end_date":    task.EndDate,
			"completed":   task.Completed,
			"reminders":   task.Reminders,
			"priority":    task.Priority,
			"tags":        task.Tags,
			"updated_at":  task.UpdatedAt,
		},
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a task by ID
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
