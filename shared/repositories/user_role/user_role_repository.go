package userrole

import (
	"context"
	"errors"
	"time"

	"github.com/atomic-blend/backend/shared/models"
	"github.com/atomic-blend/backend/shared/utils/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userRoleCollection = "user_roles"

// Interface defines methods that a UserRoleRepository must implement
type Interface interface {
	Create(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.UserRoleEntity, error)
	GetAll(ctx context.Context) ([]*models.UserRoleEntity, error)
	GetByName(ctx context.Context, name string) (*models.UserRoleEntity, error)
	Update(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindOrCreate(ctx context.Context, roleName string) (*models.UserRoleEntity, error)
	PopulateRoles(context context.Context, user *models.UserEntity) error
}

// Repository provides methods to interact with user role data in the database
type Repository struct {
	collection *mongo.Collection
}

// Ensure UserRoleRepository implements UserRoleRepositoryInterface
var _ Interface = (*Repository)(nil)

// NewUserRoleRepository creates a new instance of UserRoleRepository
func NewUserRoleRepository(database *mongo.Database) *Repository {
	if database == nil {
		database = db.Database
	}
	return &Repository{
		collection: database.Collection(userRoleCollection),
	}
}

// Create adds a new user role to the database
func (r *Repository) Create(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error) {
	// Generate an ID if not provided
	if role.ID == nil {
		id := primitive.NewObjectID()
		role.ID = &id
	}

	// Set timestamps
	now := primitive.NewDateTimeFromTime(time.Now())
	role.CreatedAt = &now
	role.UpdatedAt = &now

	// Insert into database
	_, err := r.collection.InsertOne(ctx, role)
	if err != nil {
		return nil, err
	}
	return role, nil
}

// GetByID retrieves a user role by its ID
func (r *Repository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UserRoleEntity, error) {
	var role models.UserRoleEntity
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user role not found")
		}
		return nil, err
	}
	return &role, nil
}

// GetAll retrieves all user roles
func (r *Repository) GetAll(ctx context.Context) ([]*models.UserRoleEntity, error) {
	var roles []*models.UserRoleEntity

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetByName retrieves a user role by its name
func (r *Repository) GetByName(ctx context.Context, name string) (*models.UserRoleEntity, error) {
	var role models.UserRoleEntity
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("user role not found")
		}
		return nil, err
	}
	return &role, nil
}

// Update modifies an existing user role in the database
func (r *Repository) Update(ctx context.Context, role *models.UserRoleEntity) (*models.UserRoleEntity, error) {
	if role.ID == nil {
		return nil, errors.New("role ID is required for update")
	}

	// Update timestamp
	now := primitive.NewDateTimeFromTime(time.Now())
	role.UpdatedAt = &now

	filter := bson.M{"_id": role.ID}
	update := bson.M{"$set": role}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("user role not found")
	}

	return role, nil
}

// Delete removes a user role from the database by ID
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user role not found")
	}

	return nil
}

// FindOrCreate finds a role by name or creates it if it doesn't exist
func (r *Repository) FindOrCreate(ctx context.Context, roleName string) (*models.UserRoleEntity, error) {
	// Try to find the role first
	role, err := r.GetByName(ctx, roleName)
	if err == nil {
		// Role exists, return it
		return role, nil
	}

	// Create new role if not found
	if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == "user role not found" {
		id := primitive.NewObjectID()
		now := primitive.NewDateTimeFromTime(time.Now())

		newRole := &models.UserRoleEntity{
			ID:        &id,
			Name:      roleName,
			CreatedAt: &now,
			UpdatedAt: &now,
		}

		return r.Create(ctx, newRole)
	}

	return nil, err
}

// PopulateRoles populates the roles for the given user
func (r *Repository) PopulateRoles(context context.Context, user *models.UserEntity) error {
	roles := make([]*models.UserRoleEntity, 0)
	for _, roleID := range user.RoleIds {
		role, err := r.GetByID(context, *roleID)
		if err != nil {
			return err
		}
		roles = append(roles, role)
	}
	user.Roles = roles
	return nil
}
