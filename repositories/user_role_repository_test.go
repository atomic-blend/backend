package repositories

import (
	"atomic_blend_api/models"
	"atomic_blend_api/tests/utils/inmemorymongo"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRoleDB(t *testing.T) (*UserRoleRepository, func()) {
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	mongoURI := mongoServer.URI()
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoURI)
	require.NoError(t, err)

	db := client.Database("test_db")
	repo := NewUserRoleRepository(db)

	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestUserRoleRepository_Create(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	roleName := "admin"
	role := &models.UserRoleEntity{
		Name: roleName,
	}

	created, err := repo.Create(context.Background(), role)
	assert.NoError(t, err)
	assert.NotNil(t, created.ID)
	assert.NotNil(t, created.CreatedAt)
	assert.NotNil(t, created.UpdatedAt)
	assert.Equal(t, roleName, created.Name)
}

func TestUserRoleRepository_GetByID(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	// Create a test role first
	roleName := "admin"
	role := &models.UserRoleEntity{
		Name: roleName,
	}
	created, err := repo.Create(context.Background(), role)
	require.NoError(t, err)

	// Test finding by ID
	found, err := repo.GetByID(context.Background(), *created.ID)
	assert.NoError(t, err)
	assert.Equal(t, *created.ID, *found.ID)
	assert.Equal(t, created.Name, found.Name)

	// Test with non-existent ID
	nonExistentID := primitive.NewObjectID()
	_, err = repo.GetByID(context.Background(), nonExistentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")
}

func TestUserRoleRepository_GetAll(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	// Create multiple test roles
	roles := []string{"admin", "user", "moderator"}
	for _, roleName := range roles {
		_, err := repo.Create(context.Background(), &models.UserRoleEntity{Name: roleName})
		require.NoError(t, err)
	}

	// Test getting all roles
	allRoles, err := repo.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, allRoles, len(roles))
}

func TestUserRoleRepository_GetByName(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	roleName := "admin"
	role := &models.UserRoleEntity{
		Name: roleName,
	}
	created, err := repo.Create(context.Background(), role)
	require.NoError(t, err)

	// Test finding by name
	found, err := repo.GetByName(context.Background(), roleName)
	assert.NoError(t, err)
	assert.Equal(t, *created.ID, *found.ID)
	assert.Equal(t, roleName, found.Name)

	// Test with non-existent name
	_, err = repo.GetByName(context.Background(), "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")
}

func TestUserRoleRepository_Update(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	roleName := "admin"
	role := &models.UserRoleEntity{
		Name: roleName,
	}
	created, err := repo.Create(context.Background(), role)
	require.NoError(t, err)

	// Store the original timestamp
	originalUpdatedAt := created.UpdatedAt

	// Update the role
	newName := "super_admin"
	created.Name = newName
	updated, err := repo.Update(context.Background(), created)
	assert.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.NotEqual(t, originalUpdatedAt.Time(), updated.UpdatedAt.Time())

	// Test update with non-existent ID
	nonExistentID := primitive.NewObjectID()
	created.ID = &nonExistentID
	_, err = repo.Update(context.Background(), created)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")
}

func TestUserRoleRepository_Delete(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	roleName := "admin"
	role := &models.UserRoleEntity{
		Name: roleName,
	}
	created, err := repo.Create(context.Background(), role)
	require.NoError(t, err)

	// Test deletion
	err = repo.Delete(context.Background(), *created.ID)
	assert.NoError(t, err)

	// Verify role is deleted
	_, err = repo.GetByID(context.Background(), *created.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")

	// Test delete with non-existent ID
	err = repo.Delete(context.Background(), primitive.NewObjectID())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")
}

func TestUserRoleRepository_FindOrCreate(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	roleName := "admin"

	// Test creating new role
	created, err := repo.FindOrCreate(context.Background(), roleName)
	assert.NoError(t, err)
	assert.NotNil(t, created.ID)
	assert.Equal(t, roleName, created.Name)

	// Test finding existing role
	found, err := repo.FindOrCreate(context.Background(), roleName)
	assert.NoError(t, err)
	assert.Equal(t, *created.ID, *found.ID)
	assert.Equal(t, roleName, found.Name)
}

func TestUserRoleRepository_PopulateRoles(t *testing.T) {
	repo, cleanup := setupTestRoleDB(t)
	defer cleanup()

	// Create test roles
	roles := []string{"admin", "user"}
	roleIds := make([]*primitive.ObjectID, 0)

	for _, roleName := range roles {
		role, err := repo.Create(context.Background(), &models.UserRoleEntity{Name: roleName})
		require.NoError(t, err)
		roleIds = append(roleIds, role.ID)
	}

	// Create a test user with role IDs
	user := &models.UserEntity{
		RoleIds: roleIds,
	}

	// Test populating roles
	err := repo.PopulateRoles(context.Background(), user)
	assert.NoError(t, err)
	assert.Len(t, user.Roles, len(roles))
	assert.Equal(t, "admin", user.Roles[0].Name)
	assert.Equal(t, "user", user.Roles[1].Name)

	// Test with non-existent role ID
	nonExistentID := primitive.NewObjectID()
	user.RoleIds = []*primitive.ObjectID{&nonExistentID}
	err = repo.PopulateRoles(context.Background(), user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user role not found")
}
