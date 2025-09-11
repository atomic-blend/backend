package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupDraftMailTest(t *testing.T) (DraftMailRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewDraftMailRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestCreateDraftMail(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:          &mailID,
		UserID:      userID,
		TextContent: "Test draft email content",
		Headers: map[string]string{
			"Subject": "Test Draft Email",
			"From":    "test@example.com",
		},
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)

	assert.NoError(t, err)
	assert.NotNil(t, createdDraftMail)
	assert.NotEqual(t, primitive.NilObjectID, createdDraftMail.ID)
	assert.Equal(t, models.SendStatusPending, createdDraftMail.SendStatus)
	assert.False(t, createdDraftMail.Trashed)
	assert.Nil(t, createdDraftMail.RetryCounter) // Should be nil when created
	assert.NotNil(t, createdDraftMail.CreatedAt)
	assert.NotNil(t, createdDraftMail.UpdatedAt)
}

func TestGetDraftMailByID(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:          &mailID,
		UserID:      userID,
		TextContent: "Test draft email content",
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)

	foundDraftMail, err := repository.GetByID(ctx, createdDraftMail.ID)

	assert.NoError(t, err)
	assert.NotNil(t, foundDraftMail)
	assert.Equal(t, createdDraftMail.ID, foundDraftMail.ID)
	assert.Equal(t, models.SendStatusPending, foundDraftMail.SendStatus)
}

func TestGetAllDraftMails(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()

	// Create multiple draft mails
	for i := 0; i < 5; i++ {
		mailID := primitive.NewObjectID()
		mail := &models.Mail{
			ID:     &mailID,
			UserID: userID,
		}
		draftMail := &models.SendMail{
			Mail:       mail,
			SendStatus: models.SendStatusPending,
			Trashed:    false,
		}

		_, err := repository.Create(ctx, draftMail)
		assert.NoError(t, err)
	}

	// Test pagination
	draftMails, totalCount, err := repository.GetAll(ctx, userID, 1, 3)

	assert.NoError(t, err)
	assert.Len(t, draftMails, 3)
	assert.Equal(t, int64(5), totalCount)
}

func TestGetAllDraftMailsWithoutPagination(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()

	// Create multiple draft mails
	for i := 0; i < 3; i++ {
		mailID := primitive.NewObjectID()
		mail := &models.Mail{
			ID:     &mailID,
			UserID: userID,
		}
		draftMail := &models.SendMail{
			Mail:       mail,
			SendStatus: models.SendStatusPending,
			Trashed:    false,
		}

		_, err := repository.Create(ctx, draftMail)
		assert.NoError(t, err)
	}

	// Test without pagination (page and limit <= 0)
	draftMails, totalCount, err := repository.GetAll(ctx, userID, 0, 0)

	assert.NoError(t, err)
	assert.Len(t, draftMails, 3)
	assert.Equal(t, int64(3), totalCount)
}

func TestGetAllDraftMailsWithSecondPage(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()

	// Create multiple draft mails
	for i := 0; i < 5; i++ {
		mailID := primitive.NewObjectID()
		mail := &models.Mail{
			ID:     &mailID,
			UserID: userID,
		}
		draftMail := &models.SendMail{
			Mail:       mail,
			SendStatus: models.SendStatusPending,
			Trashed:    false,
		}

		_, err := repository.Create(ctx, draftMail)
		assert.NoError(t, err)
	}

	// Test second page
	draftMails, totalCount, err := repository.GetAll(ctx, userID, 2, 3)

	assert.NoError(t, err)
	assert.Len(t, draftMails, 2) // Only 2 remaining items on second page
	assert.Equal(t, int64(5), totalCount)
}

func TestUpdateDraftMailStatus(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	update := bson.M{
		"send_status": models.SendStatusSent,
	}
	updatedDraftMail, err := repository.Update(ctx, createdDraftMail.ID, update)

	assert.NoError(t, err)
	assert.NotNil(t, updatedDraftMail)
	assert.Equal(t, models.SendStatusSent, updatedDraftMail.SendStatus)
	assert.True(t, updatedDraftMail.UpdatedAt.Time().After(createdDraftMail.UpdatedAt.Time()))
}

func TestUpdateDraftMailWithSetOperation(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	// Test with $set operation
	update := bson.M{
		"$set": bson.M{
			"send_status":    models.SendStatusFailed,
			"failure_reason": "Test failure",
		},
	}
	updatedDraftMail, err := repository.Update(ctx, createdDraftMail.ID, update)

	assert.NoError(t, err)
	assert.NotNil(t, updatedDraftMail)
	assert.Equal(t, models.SendStatusFailed, updatedDraftMail.SendStatus)
	assert.NotNil(t, updatedDraftMail.FailureReason)
	assert.Equal(t, "Test failure", *updatedDraftMail.FailureReason)
	assert.True(t, updatedDraftMail.UpdatedAt.Time().After(createdDraftMail.UpdatedAt.Time()))
}

func TestUpdateDraftMailRetryCounter(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)
	assert.Nil(t, createdDraftMail.RetryCounter) // Should start as nil

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	// Test updating retry counter and status
	retryCount := 1
	update := bson.M{
		"retry_counter": &retryCount,
		"send_status":   models.SendStatusRetry,
	}
	updatedDraftMail, err := repository.Update(ctx, createdDraftMail.ID, update)

	assert.NoError(t, err)
	assert.NotNil(t, updatedDraftMail)
	assert.Equal(t, models.SendStatusRetry, updatedDraftMail.SendStatus)
	assert.NotNil(t, updatedDraftMail.RetryCounter)
	assert.Equal(t, 1, *updatedDraftMail.RetryCounter)
	assert.True(t, updatedDraftMail.UpdatedAt.Time().After(createdDraftMail.UpdatedAt.Time()))
}

func TestTrashDraftMail(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	err = repository.Trash(ctx, createdDraftMail.ID)
	assert.NoError(t, err)

	// Verify the draft mail is marked as trashed but still exists
	foundDraftMail, err := repository.GetByID(ctx, createdDraftMail.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundDraftMail)
	assert.True(t, foundDraftMail.Trashed)
	assert.True(t, foundDraftMail.UpdatedAt.Time().After(createdDraftMail.UpdatedAt.Time()))
}

func TestDeleteDraftMail(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	draftMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdDraftMail, err := repository.Create(ctx, draftMail)
	assert.NoError(t, err)

	err = repository.Delete(ctx, createdDraftMail.ID)
	assert.NoError(t, err)

	// Verify the draft mail is actually deleted from the database
	foundDraftMail, err := repository.GetByID(ctx, createdDraftMail.ID)
	assert.NoError(t, err)
	assert.Nil(t, foundDraftMail)
}

func TestGetDraftMailByIDNotFound(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	nonExistentID := primitive.NewObjectID()
	draftMail, err := repository.GetByID(ctx, nonExistentID)

	assert.NoError(t, err)
	assert.Nil(t, draftMail)
}

func TestUpdateDraftMailNotFound(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	nonExistentID := primitive.NewObjectID()
	update := bson.M{
		"send_status": models.SendStatusSent,
	}
	updatedDraftMail, err := repository.Update(ctx, nonExistentID, update)

	assert.NoError(t, err)
	assert.Nil(t, updatedDraftMail)
}

func TestTrashDraftMailNotFound(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	nonExistentID := primitive.NewObjectID()
	err := repository.Trash(ctx, nonExistentID)

	// Trash should not return error even if document doesn't exist
	assert.NoError(t, err)
}

func TestDeleteDraftMailNotFound(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	nonExistentID := primitive.NewObjectID()
	err := repository.Delete(ctx, nonExistentID)

	// Delete should not return error even if document doesn't exist
	assert.NoError(t, err)
}

func TestDraftMailUserIsolation(t *testing.T) {
	repository, cleanup := setupDraftMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID1 := primitive.NewObjectID()
	userID2 := primitive.NewObjectID()

	// Create draft mail for user 1
	mailID1 := primitive.NewObjectID()
	mail1 := &models.Mail{
		ID:     &mailID1,
		UserID: userID1,
	}
	draftMail1 := &models.SendMail{
		Mail:       mail1,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}
	_, err := repository.Create(ctx, draftMail1)
	assert.NoError(t, err)

	// Create draft mail for user 2
	mailID2 := primitive.NewObjectID()
	mail2 := &models.Mail{
		ID:     &mailID2,
		UserID: userID2,
	}
	draftMail2 := &models.SendMail{
		Mail:       mail2,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}
	_, err = repository.Create(ctx, draftMail2)
	assert.NoError(t, err)

	// Get draft mails for user 1
	draftMails1, totalCount1, err := repository.GetAll(ctx, userID1, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, draftMails1, 1)
	assert.Equal(t, int64(1), totalCount1)

	// Get draft mails for user 2
	draftMails2, totalCount2, err := repository.GetAll(ctx, userID2, 0, 0)
	assert.NoError(t, err)
	assert.Len(t, draftMails2, 1)
	assert.Equal(t, int64(1), totalCount2)

	// Verify they are different
	assert.NotEqual(t, draftMails1[0].ID, draftMails2[0].ID)
}
