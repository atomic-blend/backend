package mocks

import (
	"context"

	"github.com/atomic-blend/backend/productivity/models"
	patchmodels "github.com/atomic-blend/backend/productivity/models/patch_models"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockNoteRepository provides a mock implementation of NoteRepository
type MockNoteRepository struct {
	mock.Mock
}

// GetAll gets all notes
func (m *MockNoteRepository) GetAll(ctx context.Context, userID *primitive.ObjectID) ([]*models.NoteEntity, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.NoteEntity), args.Error(1)
}

// GetByID gets a note by ID
func (m *MockNoteRepository) GetByID(ctx context.Context, id string) (*models.NoteEntity, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NoteEntity), args.Error(1)
}

// Create creates a new note
func (m *MockNoteRepository) Create(ctx context.Context, note *models.NoteEntity) (*models.NoteEntity, error) {
	args := m.Called(ctx, note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NoteEntity), args.Error(1)
}

// Update updates a note with the given ID
func (m *MockNoteRepository) Update(ctx context.Context, id string, note *models.NoteEntity) (*models.NoteEntity, error) {
	args := m.Called(ctx, id, note)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NoteEntity), args.Error(1)
}

// Delete deletes a note with the given ID
func (m *MockNoteRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// DeleteByUserID deletes all notes for a specific user
func (m *MockNoteRepository) DeleteByUserID(ctx context.Context, userID primitive.ObjectID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// UpdatePatch updates a note based on a patch
func (m *MockNoteRepository) UpdatePatch(ctx context.Context, patch *patchmodels.Patch) (*models.NoteEntity, error) {
	args := m.Called(ctx, patch)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.NoteEntity), args.Error(1)
}
