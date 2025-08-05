package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/mail/tests/utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupMailTest(t *testing.T) (MailRepositoryInterface, func()) {
	// Start in-memory MongoDB server
	mongoServer, err := inmemorymongo.CreateInMemoryMongoDB()
	require.NoError(t, err)

	// Connect to the in-memory MongoDB
	client, err := inmemorymongo.ConnectToInMemoryDB(mongoServer.URI())
	require.NoError(t, err)

	// Get database reference
	db := client.Database("test_db")

	repo := NewMailRepository(db)

	// Return cleanup function
	cleanup := func() {
		client.Disconnect(context.Background())
		mongoServer.Stop()
	}

	return repo, cleanup
}

func createTestMail(userID primitive.ObjectID) *models.Mail {
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Mail{
		UserID: userID,
		Headers: models.MailHeaders{
			From:      "sender@example.com",
			To:        "recipient@example.com",
			Subject:   "Test Subject",
			Date:      "2024-01-01T00:00:00Z",
			MessageID: "test-message-id-123",
		},
		TextContent: "This is a test email content",
		HTMLContent: "<html><body>This is a test email content</body></html>",
		Attachments: []models.MailAttachment{
			{
				Filename:    "test.pdf",
				ContentType: "application/pdf",
				StoragePath: "s3://bucket/test.pdf",
				StorageType: "S3",
				Size:        1024,
			},
		},
		Archived:       false,
		Trashed:        false,
		Greylisted:     false,
		Rejected:       false,
		RewriteSubject: false,
		CreatedAt:      &now,
		UpdatedAt:      &now,
	}
}

func TestMailRepository_Create(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("successful create mail", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mail := createTestMail(userID)

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)
		assert.Equal(t, userID, created.UserID)
		assert.Equal(t, mail.Headers.From, created.Headers.From)
		assert.Equal(t, mail.Headers.To, created.Headers.To)
		assert.Equal(t, mail.Headers.Subject, created.Headers.Subject)
		assert.Equal(t, mail.TextContent, created.TextContent)
		assert.Equal(t, mail.HTMLContent, created.HTMLContent)
		assert.Len(t, created.Attachments, 1)
		assert.Equal(t, "test.pdf", created.Attachments[0].Filename)
		assert.NotNil(t, created.CreatedAt)
		assert.NotNil(t, created.UpdatedAt)
	})

	t.Run("create mail with existing ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mail := createTestMail(userID)
		existingID := primitive.NewObjectID()
		mail.ID = existingID

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)
		assert.Equal(t, existingID, created.ID)
	})
}

func TestMailRepository_GetAll(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("successful get all mails for a user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create test mails for the user
		mail1 := createTestMail(userID)
		mail1.Headers.Subject = "Mail 1"

		mail2 := createTestMail(userID)
		mail2.Headers.Subject = "Mail 2"

		// Create one mail for another user
		otherUserID := primitive.NewObjectID()
		otherMail := createTestMail(otherUserID)
		otherMail.Headers.Subject = "Other User Mail"

		_, err := repo.Create(context.Background(), mail1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), mail2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), otherMail)
		require.NoError(t, err)

		// Get mails for the user
		mails, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, mails, 2)

		// Verify the mail subjects
		var subjects []string
		for _, m := range mails {
			subjects = append(subjects, m.Headers.Subject)
		}
		assert.Contains(t, subjects, "Mail 1")
		assert.Contains(t, subjects, "Mail 2")
	})

	t.Run("get all mails for user with no mails", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mails, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, mails, 0)
	})
}

func TestMailRepository_GetByID(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("successful get mail by ID", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mail := createTestMail(userID)
		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)

		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, created.UserID, found.UserID)
		assert.Equal(t, created.Headers.Subject, found.Headers.Subject)
	})

	t.Run("get mail by non-existent ID", func(t *testing.T) {
		nonExistentID := primitive.NewObjectID()
		found, err := repo.GetByID(context.Background(), nonExistentID)
		require.NoError(t, err)
		assert.Nil(t, found)
	})
}

func TestMailRepository_CreateMany(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("successful create multiple mails", func(t *testing.T) {
		// Skip this test if transactions are not supported
		// In-memory MongoDB doesn't support transactions
		t.Skip("Skipping transaction-based test - in-memory MongoDB doesn't support transactions")

		userID := primitive.NewObjectID()
		mails := []models.Mail{
			*createTestMail(userID),
			*createTestMail(userID),
			*createTestMail(userID),
		}

		// Set different subjects for each mail
		mails[0].Headers.Subject = "Batch Mail 1"
		mails[1].Headers.Subject = "Batch Mail 2"
		mails[2].Headers.Subject = "Batch Mail 3"

		success, err := repo.CreateMany(context.Background(), mails)
		require.NoError(t, err)
		assert.True(t, success)

		// Verify all mails were created
		createdMails, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, createdMails, 3)

		// Verify subjects
		var subjects []string
		for _, m := range createdMails {
			subjects = append(subjects, m.Headers.Subject)
		}
		assert.Contains(t, subjects, "Batch Mail 1")
		assert.Contains(t, subjects, "Batch Mail 2")
		assert.Contains(t, subjects, "Batch Mail 3")
	})

	t.Run("create many mails with different users", func(t *testing.T) {
		// Skip this test if transactions are not supported
		// In-memory MongoDB doesn't support transactions
		t.Skip("Skipping transaction-based test - in-memory MongoDB doesn't support transactions")

		userID1 := primitive.NewObjectID()
		userID2 := primitive.NewObjectID()

		mails := []models.Mail{
			*createTestMail(userID1),
			*createTestMail(userID2),
		}

		mails[0].Headers.Subject = "User 1 Mail"
		mails[1].Headers.Subject = "User 2 Mail"

		success, err := repo.CreateMany(context.Background(), mails)
		require.NoError(t, err)
		assert.True(t, success)

		// Verify mails for each user
		user1Mails, err := repo.GetAll(context.Background(), userID1)
		require.NoError(t, err)
		assert.Len(t, user1Mails, 1)
		assert.Equal(t, "User 1 Mail", user1Mails[0].Headers.Subject)

		user2Mails, err := repo.GetAll(context.Background(), userID2)
		require.NoError(t, err)
		assert.Len(t, user2Mails, 1)
		assert.Equal(t, "User 2 Mail", user2Mails[0].Headers.Subject)
	})

	t.Run("create empty batch", func(t *testing.T) {
		success, err := repo.CreateMany(context.Background(), []models.Mail{})
		require.NoError(t, err)
		assert.True(t, success)
	})
}

func TestMailRepository_Integration(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("full mail lifecycle", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create a mail
		mail := createTestMail(userID)
		mail.Headers.Subject = "Integration Test Mail"
		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)

		// Get the mail by ID
		found, err := repo.GetByID(context.Background(), created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Integration Test Mail", found.Headers.Subject)

		// Get all mails for the user
		allMails, err := repo.GetAll(context.Background(), userID)
		require.NoError(t, err)
		assert.Len(t, allMails, 1)
		assert.Equal(t, created.ID, allMails[0].ID)
	})
}
