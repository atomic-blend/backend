package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/test_utils/inmemorymongo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
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

// Helper function to convert headers from MongoDB format to map[string]interface{}
func convertHeadersToMap(headers interface{}) (map[string]interface{}, error) {
	if headers == nil {
		return nil, nil
	}

	// If it's already a map, return it
	if headerMap, ok := headers.(map[string]interface{}); ok {
		return headerMap, nil
	}

	// If it's a primitive.D (BSON document), convert it
	if headerDoc, ok := headers.(primitive.D); ok {
		result := make(map[string]interface{})
		for _, elem := range headerDoc {
			result[elem.Key] = elem.Value
		}
		return result, nil
	}

	// Try to marshal and unmarshal to convert any BSON type
	bsonData, err := bson.Marshal(headers)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = bson.Unmarshal(bsonData, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createTestMail(userID primitive.ObjectID) *models.Mail {
	now := primitive.NewDateTimeFromTime(time.Now())
	return &models.Mail{
		UserID: userID,
		Headers: map[string]string{
			"From":       "sender@example.com",
			"To":         "recipient@example.com",
			"Subject":    "Test Subject",
			"Date":       "2024-01-01T00:00:00Z",
			"Message-Id": "test-message-id-123",
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
		Archived:       nil,
		Trashed:        nil,
		Greylisted:     nil,
		Rejected:       nil,
		RewriteSubject: nil,
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

		// Compare headers using the helper function
		originalHeaders, err := convertHeadersToMap(mail.Headers)
		require.NoError(t, err)
		createdHeaders, err := convertHeadersToMap(created.Headers)
		require.NoError(t, err)
		assert.Equal(t, originalHeaders, createdHeaders)

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
		mail.ID = &existingID

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)
		assert.Equal(t, &existingID, created.ID)
	})
}

func TestMailRepository_GetAll(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("successful get all mails for a user", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create test mails for the user
		mail1 := createTestMail(userID)
		if headers1, ok := mail1.Headers.(map[string]string); ok {
			headers1["Subject"] = "Mail 1"
		}

		mail2 := createTestMail(userID)
		if headers2, ok := mail2.Headers.(map[string]string); ok {
			headers2["Subject"] = "Mail 2"
		}

		// Create one mail for another user
		otherUserID := primitive.NewObjectID()
		otherMail := createTestMail(otherUserID)
		if headersOther, ok := otherMail.Headers.(map[string]string); ok {
			headersOther["Subject"] = "Other User Mail"
		}

		_, err := repo.Create(context.Background(), mail1)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), mail2)
		require.NoError(t, err)

		_, err = repo.Create(context.Background(), otherMail)
		require.NoError(t, err)

		// Get mails for the user
		mails, total, err := repo.GetAll(context.Background(), userID, 0, 0)
		require.NoError(t, err)
		assert.Len(t, mails, 2)
		assert.Equal(t, int64(2), total)

		// Verify the mail subjects
		var subjects []string
		for _, m := range mails {
			headers, err := convertHeadersToMap(m.Headers)
			require.NoError(t, err)
			if headers != nil {
				if subject, exists := headers["Subject"]; exists {
					if subjectStr, ok := subject.(string); ok {
						subjects = append(subjects, subjectStr)
					}
				}
			}
		}
		assert.Contains(t, subjects, "Mail 1")
		assert.Contains(t, subjects, "Mail 2")
	})

	t.Run("get all mails for user with no mails", func(t *testing.T) {
		userID := primitive.NewObjectID()
		mails, total, err := repo.GetAll(context.Background(), userID, 0, 0)
		require.NoError(t, err)
		assert.Len(t, mails, 0)
		assert.Equal(t, int64(0), total)
	})

	// Add a paginated retrieval test
	t.Run("get mails paginated", func(t *testing.T) {
		userID := primitive.NewObjectID()
		// Create 5 mails
		for i := 0; i < 5; i++ {
			mail := createTestMail(userID)
			if headers, ok := mail.Headers.(map[string]string); ok {
				headers["Subject"] = fmt.Sprintf("Mail %d", i+1)
			}
			_, err := repo.Create(context.Background(), mail)
			require.NoError(t, err)
		}
		// Get first 2 mails (page 1, limit 2)
		mails, total, err := repo.GetAll(context.Background(), userID, 1, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 2)
		assert.Equal(t, int64(5), total)
		// Get next 2 mails (page 2, limit 2)
		mails, total, err = repo.GetAll(context.Background(), userID, 2, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 2)
		assert.Equal(t, int64(5), total)
		// Get last mail (page 3, limit 2)
		mails, total, err = repo.GetAll(context.Background(), userID, 3, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 1)
		assert.Equal(t, int64(5), total)
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

		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)
		assert.Equal(t, created.UserID, found.UserID)

		// Compare headers using the helper function
		createdHeaders, err := convertHeadersToMap(created.Headers)
		require.NoError(t, err)
		foundHeaders, err := convertHeadersToMap(found.Headers)
		require.NoError(t, err)
		assert.Equal(t, createdHeaders, foundHeaders)
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
		if headers0, ok := mails[0].Headers.(map[string]string); ok {
			headers0["Subject"] = "Batch Mail 1"
		}
		if headers1, ok := mails[1].Headers.(map[string]string); ok {
			headers1["Subject"] = "Batch Mail 2"
		}
		if headers2, ok := mails[2].Headers.(map[string]string); ok {
			headers2["Subject"] = "Batch Mail 3"
		}

		success, err := repo.CreateMany(context.Background(), mails)
		require.NoError(t, err)
		assert.True(t, success)

		// Verify all mails were created
		createdMails, total, err := repo.GetAll(context.Background(), userID, 0, 0)
		require.NoError(t, err)
		assert.Len(t, createdMails, 3)
		assert.Equal(t, int64(3), total)

		// Verify subjects
		var subjects []string
		for _, m := range createdMails {
			headers, err := convertHeadersToMap(m.Headers)
			require.NoError(t, err)
			if headers != nil {
				if subject, exists := headers["Subject"]; exists {
					if subjectStr, ok := subject.(string); ok {
						subjects = append(subjects, subjectStr)
					}
				}
			}
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

		if headers0, ok := mails[0].Headers.(map[string]string); ok {
			headers0["Subject"] = "User 1 Mail"
		}
		if headers1, ok := mails[1].Headers.(map[string]string); ok {
			headers1["Subject"] = "User 2 Mail"
		}

		success, err := repo.CreateMany(context.Background(), mails)
		require.NoError(t, err)
		assert.True(t, success)

		// Verify mails for each user
		user1Mails, total, err := repo.GetAll(context.Background(), userID1, 0, 0)
		require.NoError(t, err)
		assert.Len(t, user1Mails, 1)
		assert.Equal(t, int64(1), total)

		user2Mails, total, err := repo.GetAll(context.Background(), userID2, 0, 0)
		require.NoError(t, err)
		assert.Len(t, user2Mails, 1)
		assert.Equal(t, int64(1), total)
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
		if headers, ok := mail.Headers.(map[string]string); ok {
			headers["Subject"] = "Integration Test Mail"
		}
		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)
		assert.NotNil(t, created.ID)

		// Get the mail by ID
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, *created.ID, *found.ID)

		// Verify the subject using helper function
		headers, err := convertHeadersToMap(found.Headers)
		require.NoError(t, err)
		if headers != nil {
			if subject, exists := headers["Subject"]; exists {
				if subjectStr, ok := subject.(string); ok {
					assert.Equal(t, "Integration Test Mail", subjectStr)
				}
			}
		}

		// Get all mails for the user
		allMails, total, err := repo.GetAll(context.Background(), userID, 0, 0)
		require.NoError(t, err)
		assert.Len(t, allMails, 1)
		assert.Equal(t, int64(1), total)
	})
}
