package repositories

import (
	"atomic-blend/backend/productivity/models"
	patchmodels "atomic-blend/backend/productivity/models/patch_models"
	keyconverter "atomic-blend/backend/productivity/utils/key_converter"
	"context"
	"errors"
	"time"

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
	UpdatePatch(ctx context.Context, patch *patchmodels.Patch) (*models.TaskEntity, error)
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
		"_id":         objID,
		"title":       task.Title,
		"user":        task.User,
		"description": task.Description,
		"start_date":  task.StartDate,
		"end_date":    task.EndDate,
		"completed":   task.Completed,
		"reminders":   task.Reminders,
		"priority":    task.Priority,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
	})

	if err != nil {
		return nil, err
	}

	return task, nil
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
			"folder_id":   task.FolderID,
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

// UpdatePatch applies a patch to a task based on the provided patch model
func (r *TaskRepository) UpdatePatch(ctx context.Context, patch *patchmodels.Patch) (*models.TaskEntity, error) {
	if patch.Action != "update" {
		return nil, errors.New("only update action is supported")
	}

	if patch.ItemType != patchmodels.ItemTypeTask {
		return nil, errors.New("item type not supported")
	}

	updatePayload := bson.M{}
	for _, change := range patch.Changes {
		//convert Key from camelCase to snake_case
		key := keyconverter.ToSnakeCase(change.Key)
		print("Key: ", change.Key, "\n")
		print("converted: ", keyconverter.ToSnakeCase(change.Key), "\n")
		value := change.Value
		if isDateTimeField(key) {
			if dateValue, err := convertToDateTime(change.Value, isDateTimePointerField(key)); err == nil {
				value = dateValue
			} else {
				return nil, errors.New("invalid date format for field: " + key)
			}
		} else if isBooleanField(key) {
			if boolValue, err := convertToBoolean(change.Value, isBooleanPointerField(key)); err == nil {
				value = boolValue
			} else {
				return nil, errors.New("invalid boolean format for field: " + key)
			}
		}

		updatePayload[key] = value
	}

	updatePayload["updated_at"] = primitive.NewDateTimeFromTime(time.Now())

	// Perform the update operation
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": patch.ItemID}, bson.M{"$set": updatePayload})
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, patch.ItemID.Hex())
}

// Helper function to check if a field is a date/time field
func isDateTimeField(fieldName string) bool {
	dateTimeFields := []string{"start_date", "end_date", "created_at", "updated_at"}
	for _, field := range dateTimeFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to check if a field should be a pointer to DateTime
func isDateTimePointerField(fieldName string) bool {
	pointerFields := []string{"start_date", "end_date"}
	for _, field := range pointerFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to check if a field is a boolean field
func isBooleanField(fieldName string) bool {
	booleanFields := []string{"completed"}
	for _, field := range booleanFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to check if a field should be a pointer to bool
func isBooleanPointerField(fieldName string) bool {
	pointerFields := []string{"completed"}
	for _, field := range pointerFields {
		if fieldName == field {
			return true
		}
	}
	return false
}

// Helper function to convert various date formats to primitive.DateTime
func convertToDateTime(value interface{}, isPointer bool) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	var parsedTime time.Time
	var err error

	switch v := value.(type) {
	case string:
		// Try parsing different date formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02T15:04:05.000Z",
			"2006-01-02 15:04:05",
		}

		for _, format := range formats {
			if parsedTime, err = time.Parse(format, v); err == nil {
				break
			}
		}
		if err != nil {
			return nil, errors.New("unable to parse date string")
		}
	case int64:
		// Unix timestamp in milliseconds
		parsedTime = time.Unix(0, v*int64(time.Millisecond))
	case float64:
		// Unix timestamp in seconds (JavaScript timestamps are often float64)
		parsedTime = time.Unix(int64(v), 0)
	default:
		return nil, errors.New("unsupported date format")
	}

	dt := primitive.NewDateTimeFromTime(parsedTime)

	if isPointer {
		return &dt, nil
	}
	return dt, nil
}

// Helper function to convert various boolean formats to bool
func convertToBoolean(value interface{}, isPointer bool) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	var boolValue bool

	switch v := value.(type) {
	case bool:
		boolValue = v
	case string:
		switch v {
		case "true", "True", "TRUE", "1":
			boolValue = true
		case "false", "False", "FALSE", "0":
			boolValue = false
		default:
			return nil, errors.New("unable to parse boolean string")
		}
	case int:
		boolValue = v != 0
	case int64:
		boolValue = v != 0
	case float64:
		boolValue = v != 0
	default:
		return nil, errors.New("unsupported boolean format")
	}

	if isPointer {
		return &boolValue, nil
	}
	return boolValue, nil
}
