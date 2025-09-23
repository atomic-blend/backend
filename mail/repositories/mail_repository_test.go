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
		mail2 := createTestMail(userID)

		// Ensure deterministic CreatedAt ordering: mail2 is newer than mail1
		base := time.Now()
		createdAt1 := primitive.NewDateTimeFromTime(base.Add(-1 * time.Millisecond))
		createdAt2 := primitive.NewDateTimeFromTime(base)
		mail1.CreatedAt = &createdAt1
		mail1.UpdatedAt = &createdAt1
		mail2.CreatedAt = &createdAt2
		mail2.UpdatedAt = &createdAt2

		if headers1, ok := mail1.Headers.(map[string]string); ok {
			headers1["Subject"] = "Mail 1"
		}

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

		// Ensure mails are sorted from most recent to oldest
		if len(mails) >= 2 {
			headers0, err := convertHeadersToMap(mails[0].Headers)
			require.NoError(t, err)
			headers1, err := convertHeadersToMap(mails[1].Headers)
			require.NoError(t, err)
			if headers0 != nil && headers1 != nil {
				if s0, ok := headers0["Subject"].(string); ok {
					if s1, ok := headers1["Subject"].(string); ok {
						// mail2 was created after mail1, so it should come first
						assert.Equal(t, "Mail 2", s0)
						assert.Equal(t, "Mail 1", s1)
					}
				}
			}
		}
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
		// Create 5 mails with deterministic CreatedAt timestamps
		baseTime := time.Now()
		for i := 0; i < 5; i++ {
			mail := createTestMail(userID)
			// Ensure CreatedAt increases with i so Mail 1 is oldest and Mail 5 is newest
			createdAt := primitive.NewDateTimeFromTime(baseTime.Add(time.Duration(i) * time.Millisecond))
			mail.CreatedAt = &createdAt
			mail.UpdatedAt = &createdAt
			if headers, ok := mail.Headers.(map[string]string); ok {
				headers["Subject"] = fmt.Sprintf("Mail %d", i+1)
			}
			_, err := repo.Create(context.Background(), mail)
			require.NoError(t, err)
		}
		// Get first 2 mails (page 1, limit 2) -> should be Mail 5, Mail 4
		mails, total, err := repo.GetAll(context.Background(), userID, 1, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 2)
		assert.Equal(t, int64(5), total)
		if len(mails) >= 2 {
			headers0, err := convertHeadersToMap(mails[0].Headers)
			require.NoError(t, err)
			headers1, err := convertHeadersToMap(mails[1].Headers)
			require.NoError(t, err)
			if headers0 != nil && headers1 != nil {
				if s0, ok := headers0["Subject"].(string); ok {
					if s1, ok := headers1["Subject"].(string); ok {
						assert.Equal(t, "Mail 5", s0)
						assert.Equal(t, "Mail 4", s1)
					}
				}
			}
		}
		// Get next 2 mails (page 2, limit 2) -> should be Mail 3, Mail 2
		mails, total, err = repo.GetAll(context.Background(), userID, 2, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 2)
		assert.Equal(t, int64(5), total)
		if len(mails) >= 2 {
			headers0, err := convertHeadersToMap(mails[0].Headers)
			require.NoError(t, err)
			headers1, err := convertHeadersToMap(mails[1].Headers)
			require.NoError(t, err)
			if headers0 != nil && headers1 != nil {
				if s0, ok := headers0["Subject"].(string); ok {
					if s1, ok := headers1["Subject"].(string); ok {
						assert.Equal(t, "Mail 3", s0)
						assert.Equal(t, "Mail 2", s1)
					}
				}
			}
		}
		// Get last mail (page 3, limit 2) -> should be Mail 1
		mails, total, err = repo.GetAll(context.Background(), userID, 3, 2)
		require.NoError(t, err)
		assert.Len(t, mails, 1)
		assert.Equal(t, int64(5), total)
		if len(mails) == 1 {
			headers0, err := convertHeadersToMap(mails[0].Headers)
			require.NoError(t, err)
			if headers0 != nil {
				if s0, ok := headers0["Subject"].(string); ok {
					assert.Equal(t, "Mail 1", s0)
				}
			}
		}
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

func TestMailRepository_CleanupTrash(t *testing.T) {
	repo, cleanup := setupMailTest(t)
	defer cleanup()

	t.Run("cleanup old trashed mails", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now()

		// Create mails with different trashed_at dates
		testCases := []struct {
			name        string
			trashed     bool
			trashedAt   *primitive.DateTime
			shouldExist bool
		}{
			{
				name:        "old trashed mail (35 days ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35)); return &dt }(),
				shouldExist: false, // Should be deleted
			},
			{
				name:        "recent trashed mail (10 days ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -10)); return &dt }(),
				shouldExist: true, // Should remain
			},
			{
				name:        "old non-trashed mail (35 days ago)",
				trashed:     false,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35)); return &dt }(),
				shouldExist: true, // Should remain (not trashed)
			},
			{
				name:        "trashed mail with nil trashed_at",
				trashed:     true,
				trashedAt:   nil,
				shouldExist: true, // Should remain (no trashed_at date)
			},
			{
				name:        "non-trashed mail with old trashed_at",
				trashed:     false,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35)); return &dt }(),
				shouldExist: true, // Should remain (not trashed)
			},
		}

		var createdMails []*models.Mail
		for i, tc := range testCases {
			mail := createTestMail(userID)
			if headers, ok := mail.Headers.(map[string]string); ok {
				headers["Subject"] = tc.name
			}
			mail.Trashed = &tc.trashed
			mail.TrashedAt = tc.trashedAt

			created, err := repo.Create(context.Background(), mail)
			require.NoError(t, err)
			createdMails = append(createdMails, created)

			// Verify mail was created
			found, err := repo.GetByID(context.Background(), *created.ID)
			require.NoError(t, err)
			assert.NotNil(t, found, "Mail %d (%s) should exist before cleanup", i, tc.name)
		}

		// Run cleanup (default 30 days)
		err := repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)

		// Verify results
		for i, tc := range testCases {
			found, err := repo.GetByID(context.Background(), *createdMails[i].ID)
			require.NoError(t, err)

			if tc.shouldExist {
				assert.NotNil(t, found, "Mail %d (%s) should still exist after cleanup", i, tc.name)
			} else {
				assert.Nil(t, found, "Mail %d (%s) should be deleted after cleanup", i, tc.name)
			}
		}
	})

	t.Run("cleanup with no trashed mails", func(t *testing.T) {
		userID := primitive.NewObjectID()

		// Create a non-trashed mail
		mail := createTestMail(userID)
		trashed := false
		mail.Trashed = &trashed

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)

		// Run cleanup
		err = repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)

		// Verify mail still exists
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found)
	})

	t.Run("cleanup with no mails at all", func(t *testing.T) {
		// Run cleanup on empty collection
		err := repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)
	})

	t.Run("cleanup with exactly 30 days old trashed mail", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now()

		// Create a mail trashed exactly 30 days ago (should be deleted)
		mail := createTestMail(userID)
		trashed := true
		mail.Trashed = &trashed
		trashedAt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -30))
		mail.TrashedAt = &trashedAt

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)

		// Run cleanup
		err = repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)

		// Verify mail is deleted (30 days is the cutoff)
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.Nil(t, found, "Mail trashed exactly 30 days ago should be deleted")
	})

	t.Run("cleanup with 29 days old trashed mail", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now()

		// Create a mail trashed 29 days ago (should remain)
		mail := createTestMail(userID)
		trashed := true
		mail.Trashed = &trashed
		trashedAt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -29))
		mail.TrashedAt = &trashedAt

		created, err := repo.Create(context.Background(), mail)
		require.NoError(t, err)

		// Run cleanup
		err = repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)

		// Verify mail still exists
		found, err := repo.GetByID(context.Background(), *created.ID)
		require.NoError(t, err)
		assert.NotNil(t, found, "Mail trashed 29 days ago should remain")
	})

	t.Run("cleanup with multiple users - global cleanup", func(t *testing.T) {
		userID1 := primitive.NewObjectID()
		userID2 := primitive.NewObjectID()
		now := time.Now()

		// Create old trashed mails for both users
		mail1 := createTestMail(userID1)
		trashed1 := true
		mail1.Trashed = &trashed1
		trashedAt1 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35))
		mail1.TrashedAt = &trashedAt1

		mail2 := createTestMail(userID2)
		trashed2 := true
		mail2.Trashed = &trashed2
		trashedAt2 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -40))
		mail2.TrashedAt = &trashedAt2

		created1, err := repo.Create(context.Background(), mail1)
		require.NoError(t, err)
		created2, err := repo.Create(context.Background(), mail2)
		require.NoError(t, err)

		// Run global cleanup (userID = nil)
		err = repo.CleanupTrash(context.Background(), nil, nil)
		require.NoError(t, err)

		// Verify both mails are deleted
		found1, err := repo.GetByID(context.Background(), *created1.ID)
		require.NoError(t, err)
		assert.Nil(t, found1, "User 1's old trashed mail should be deleted")

		found2, err := repo.GetByID(context.Background(), *created2.ID)
		require.NoError(t, err)
		assert.Nil(t, found2, "User 2's old trashed mail should be deleted")
	})

	t.Run("cleanup with user-specific filtering", func(t *testing.T) {
		userID1 := primitive.NewObjectID()
		userID2 := primitive.NewObjectID()
		now := time.Now()

		// Create old trashed mails for both users
		mail1 := createTestMail(userID1)
		trashed1 := true
		mail1.Trashed = &trashed1
		trashedAt1 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35))
		mail1.TrashedAt = &trashedAt1

		mail2 := createTestMail(userID2)
		trashed2 := true
		mail2.Trashed = &trashed2
		trashedAt2 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -40))
		mail2.TrashedAt = &trashedAt2

		created1, err := repo.Create(context.Background(), mail1)
		require.NoError(t, err)
		created2, err := repo.Create(context.Background(), mail2)
		require.NoError(t, err)

		// Run cleanup only for userID1
		err = repo.CleanupTrash(context.Background(), &userID1, nil)
		require.NoError(t, err)

		// Verify only userID1's mail is deleted
		found1, err := repo.GetByID(context.Background(), *created1.ID)
		require.NoError(t, err)
		assert.Nil(t, found1, "User 1's old trashed mail should be deleted")

		// Verify userID2's mail still exists
		found2, err := repo.GetByID(context.Background(), *created2.ID)
		require.NoError(t, err)
		assert.NotNil(t, found2, "User 2's old trashed mail should still exist")
	})

	t.Run("cleanup with user-specific filtering - no matching user mails", func(t *testing.T) {
		userID1 := primitive.NewObjectID()
		userID2 := primitive.NewObjectID()
		now := time.Now()

		// Create old trashed mail only for userID2
		mail2 := createTestMail(userID2)
		trashed2 := true
		mail2.Trashed = &trashed2
		trashedAt2 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35))
		mail2.TrashedAt = &trashedAt2

		created2, err := repo.Create(context.Background(), mail2)
		require.NoError(t, err)

		// Run cleanup only for userID1 (who has no mails)
		err = repo.CleanupTrash(context.Background(), &userID1, nil)
		require.NoError(t, err)

		// Verify userID2's mail still exists (should not be affected)
		found2, err := repo.GetByID(context.Background(), *created2.ID)
		require.NoError(t, err)
		assert.NotNil(t, found2, "User 2's mail should still exist when cleaning user 1")
	})

	t.Run("cleanup with user-specific filtering - mixed scenarios", func(t *testing.T) {
		userID1 := primitive.NewObjectID()
		userID2 := primitive.NewObjectID()
		now := time.Now()

		// Create various mails for userID1
		mail1OldTrashed := createTestMail(userID1)
		trashed1 := true
		mail1OldTrashed.Trashed = &trashed1
		trashedAt1 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35))
		mail1OldTrashed.TrashedAt = &trashedAt1

		mail1RecentTrashed := createTestMail(userID1)
		mail1RecentTrashed.Trashed = &trashed1
		trashedAtRecent := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -10))
		mail1RecentTrashed.TrashedAt = &trashedAtRecent

		mail1NotTrashed := createTestMail(userID1)
		notTrashed := false
		mail1NotTrashed.Trashed = &notTrashed

		// Create old trashed mail for userID2
		mail2OldTrashed := createTestMail(userID2)
		mail2OldTrashed.Trashed = &trashed1
		trashedAt2 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -40))
		mail2OldTrashed.TrashedAt = &trashedAt2

		created1Old, err := repo.Create(context.Background(), mail1OldTrashed)
		require.NoError(t, err)
		created1Recent, err := repo.Create(context.Background(), mail1RecentTrashed)
		require.NoError(t, err)
		created1Not, err := repo.Create(context.Background(), mail1NotTrashed)
		require.NoError(t, err)
		created2Old, err := repo.Create(context.Background(), mail2OldTrashed)
		require.NoError(t, err)

		// Run cleanup only for userID1
		err = repo.CleanupTrash(context.Background(), &userID1, nil)
		require.NoError(t, err)

		// Verify userID1's old trashed mail is deleted
		found1Old, err := repo.GetByID(context.Background(), *created1Old.ID)
		require.NoError(t, err)
		assert.Nil(t, found1Old, "User 1's old trashed mail should be deleted")

		// Verify userID1's recent trashed mail still exists
		found1Recent, err := repo.GetByID(context.Background(), *created1Recent.ID)
		require.NoError(t, err)
		assert.NotNil(t, found1Recent, "User 1's recent trashed mail should still exist")

		// Verify userID1's non-trashed mail still exists
		found1Not, err := repo.GetByID(context.Background(), *created1Not.ID)
		require.NoError(t, err)
		assert.NotNil(t, found1Not, "User 1's non-trashed mail should still exist")

		// Verify userID2's old trashed mail still exists (not affected by userID1 cleanup)
		found2Old, err := repo.GetByID(context.Background(), *created2Old.ID)
		require.NoError(t, err)
		assert.NotNil(t, found2Old, "User 2's old trashed mail should still exist")
	})

	t.Run("cleanup with user-specific filtering - non-existent user", func(t *testing.T) {
		userID1 := primitive.NewObjectID()
		nonExistentUserID := primitive.NewObjectID()
		now := time.Now()

		// Create old trashed mail for userID1
		mail1 := createTestMail(userID1)
		trashed1 := true
		mail1.Trashed = &trashed1
		trashedAt1 := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -35))
		mail1.TrashedAt = &trashedAt1

		created1, err := repo.Create(context.Background(), mail1)
		require.NoError(t, err)

		// Run cleanup for non-existent user
		err = repo.CleanupTrash(context.Background(), &nonExistentUserID, nil)
		require.NoError(t, err)

		// Verify userID1's mail still exists (not affected)
		found1, err := repo.GetByID(context.Background(), *created1.ID)
		require.NoError(t, err)
		assert.NotNil(t, found1, "User 1's mail should still exist when cleaning non-existent user")
	})

	t.Run("cleanup with custom days parameter - 7 days", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now()

		// Create mails with different trashed_at dates
		testCases := []struct {
			name        string
			trashed     bool
			trashedAt   *primitive.DateTime
			shouldExist bool
		}{
			{
				name:        "old trashed mail (10 days ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -10)); return &dt }(),
				shouldExist: false, // Should be deleted (older than 7 days)
			},
			{
				name:        "recent trashed mail (5 days ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -5)); return &dt }(),
				shouldExist: true, // Should remain (newer than 7 days)
			},
			{
				name:        "exactly 7 days old trashed mail",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -7)); return &dt }(),
				shouldExist: false, // Should be deleted (exactly 7 days)
			},
		}

		var createdMails []*models.Mail
		for _, tc := range testCases {
			mail := createTestMail(userID)
			if headers, ok := mail.Headers.(map[string]string); ok {
				headers["Subject"] = tc.name
			}
			mail.Trashed = &tc.trashed
			mail.TrashedAt = tc.trashedAt

			created, err := repo.Create(context.Background(), mail)
			require.NoError(t, err)
			createdMails = append(createdMails, created)
		}

		// Run cleanup with 7 days parameter
		days := 7
		err := repo.CleanupTrash(context.Background(), &userID, &days)
		require.NoError(t, err)

		// Verify results
		for i, tc := range testCases {
			found, err := repo.GetByID(context.Background(), *createdMails[i].ID)
			require.NoError(t, err)

			if tc.shouldExist {
				assert.NotNil(t, found, "Mail %d (%s) should still exist after cleanup", i, tc.name)
			} else {
				assert.Nil(t, found, "Mail %d (%s) should be deleted after cleanup", i, tc.name)
			}
		}
	})

	t.Run("cleanup with days parameter -1 (delete all trashed)", func(t *testing.T) {
		userID := primitive.NewObjectID()
		now := time.Now()

		// Create mails with different trashed_at dates
		testCases := []struct {
			name        string
			trashed     bool
			trashedAt   *primitive.DateTime
			shouldExist bool
		}{
			{
				name:        "recent trashed mail (1 day ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -1)); return &dt }(),
				shouldExist: false, // Should be deleted (all trashed mails)
			},
			{
				name:        "very recent trashed mail (1 hour ago)",
				trashed:     true,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.Add(-1 * time.Hour)); return &dt }(),
				shouldExist: false, // Should be deleted (all trashed mails)
			},
			{
				name:        "non-trashed mail",
				trashed:     false,
				trashedAt:   func() *primitive.DateTime { dt := primitive.NewDateTimeFromTime(now.AddDate(0, 0, -1)); return &dt }(),
				shouldExist: true, // Should remain (not trashed)
			},
		}

		var createdMails []*models.Mail
		for _, tc := range testCases {
			mail := createTestMail(userID)
			if headers, ok := mail.Headers.(map[string]string); ok {
				headers["Subject"] = tc.name
			}
			mail.Trashed = &tc.trashed
			mail.TrashedAt = tc.trashedAt

			created, err := repo.Create(context.Background(), mail)
			require.NoError(t, err)
			createdMails = append(createdMails, created)
		}

		// Run cleanup with days = -1 (delete all trashed)
		days := -1
		err := repo.CleanupTrash(context.Background(), &userID, &days)
		require.NoError(t, err)

		// Verify results
		for i, tc := range testCases {
			found, err := repo.GetByID(context.Background(), *createdMails[i].ID)
			require.NoError(t, err)

			if tc.shouldExist {
				assert.NotNil(t, found, "Mail %d (%s) should still exist after cleanup", i, tc.name)
			} else {
				assert.Nil(t, found, "Mail %d (%s) should be deleted after cleanup", i, tc.name)
			}
		}
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
