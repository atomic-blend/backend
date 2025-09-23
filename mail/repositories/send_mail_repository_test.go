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

func setupSendMailTest(t *testing.T) (SendMailRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewSendMailRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func TestCreateSendMail(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:          &mailID,
		UserID:      userID,
		TextContent: "Test email content",
		Headers: map[string]string{
			"Subject": "Test Email",
			"From":    "test@example.com",
		},
	}

	sendMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdSendMail, err := repository.Create(ctx, sendMail)

	assert.NoError(t, err)
	assert.NotNil(t, createdSendMail)
	assert.NotEqual(t, primitive.NilObjectID, createdSendMail.ID)
	assert.Equal(t, models.SendStatusPending, createdSendMail.SendStatus)
	assert.False(t, createdSendMail.Trashed)
	assert.Nil(t, createdSendMail.RetryCounter) // Should be nil when created
	assert.NotNil(t, createdSendMail.CreatedAt)
	assert.NotNil(t, createdSendMail.UpdatedAt)
}

func TestGetSendMailByID(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:          &mailID,
		UserID:      userID,
		TextContent: "Test email content",
	}

	sendMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdSendMail, err := repository.Create(ctx, sendMail)
	assert.NoError(t, err)

	foundSendMail, err := repository.GetByID(ctx, createdSendMail.ID)

	assert.NoError(t, err)
	assert.NotNil(t, foundSendMail)
	assert.Equal(t, createdSendMail.ID, foundSendMail.ID)
	assert.Equal(t, models.SendStatusPending, foundSendMail.SendStatus)
}

func TestGetAllSendMails(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()

	// Create multiple send mails
	for i := 0; i < 5; i++ {
		mailID := primitive.NewObjectID()
		mail := &models.Mail{
			ID:     &mailID,
			UserID: userID,
		}
		sendMail := &models.SendMail{
			Mail:       mail,
			SendStatus: models.SendStatusPending,
			Trashed:    false,
		}

		_, err := repository.Create(ctx, sendMail)
		assert.NoError(t, err)
	}

	// Test pagination
	sendMails, totalCount, err := repository.GetAll(ctx, userID, 1, 3)

	assert.NoError(t, err)
	assert.Len(t, sendMails, 3)
	assert.Equal(t, int64(5), totalCount)
}

func TestUpdateSendMailStatus(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	sendMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdSendMail, err := repository.Create(ctx, sendMail)
	assert.NoError(t, err)

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	update := bson.M{
		"send_status": models.SendStatusSent,
	}
	updatedSendMail, err := repository.Update(ctx, createdSendMail.ID, update)

	assert.NoError(t, err)
	assert.NotNil(t, updatedSendMail)
	assert.Equal(t, models.SendStatusSent, updatedSendMail.SendStatus)
	assert.True(t, updatedSendMail.UpdatedAt.Time().After(createdSendMail.UpdatedAt.Time()))
}

func TestUpdateSendMailRetryCounter(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	sendMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdSendMail, err := repository.Create(ctx, sendMail)
	assert.NoError(t, err)
	assert.Nil(t, createdSendMail.RetryCounter) // Should start as nil

	// Wait a moment to ensure different timestamp
	time.Sleep(time.Millisecond * 10)

	// Test updating retry counter and status
	retryCount := 1
	update := bson.M{
		"retry_counter": &retryCount,
		"send_status":   models.SendStatusRetry,
	}
	updatedSendMail, err := repository.Update(ctx, createdSendMail.ID, update)

	assert.NoError(t, err)
	assert.NotNil(t, updatedSendMail)
	assert.Equal(t, models.SendStatusRetry, updatedSendMail.SendStatus)
	assert.NotNil(t, updatedSendMail.RetryCounter)
	assert.Equal(t, 1, *updatedSendMail.RetryCounter)
	assert.True(t, updatedSendMail.UpdatedAt.Time().After(createdSendMail.UpdatedAt.Time()))
}

func TestDeleteSendMail(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	userID := primitive.NewObjectID()
	mailID := primitive.NewObjectID()
	mail := &models.Mail{
		ID:     &mailID,
		UserID: userID,
	}

	sendMail := &models.SendMail{
		Mail:       mail,
		SendStatus: models.SendStatusPending,
		Trashed:    false,
	}

	createdSendMail, err := repository.Create(ctx, sendMail)
	assert.NoError(t, err)

	err = repository.Delete(ctx, createdSendMail.ID)
	assert.NoError(t, err)

	// Verify the send mail is marked as trashed
	foundSendMail, err := repository.GetByID(ctx, createdSendMail.ID)
	assert.NoError(t, err)
	assert.NotNil(t, foundSendMail)
	assert.True(t, foundSendMail.Trashed)
}

func TestGetSendMailByIDNotFound(t *testing.T) {
	repository, cleanup := setupSendMailTest(t)
	defer cleanup()

	ctx := context.Background()

	nonExistentID := primitive.NewObjectID()
	sendMail, err := repository.GetByID(ctx, nonExistentID)

	assert.NoError(t, err)
	assert.Nil(t, sendMail)
}
