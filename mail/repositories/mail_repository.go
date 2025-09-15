package repositories

import (
	"context"
	"time"

	"github.com/atomic-blend/backend/mail/models"
	"github.com/atomic-blend/backend/shared/utils/db"

	bson "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mailCollection = "mails"

// MailRepositoryInterface defines the interface for mail repository operations
type MailRepositoryInterface interface {
	// GetAll retrieves mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all mails and total count.
	GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.Mail, int64, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*models.Mail, error)
	Create(ctx context.Context, mail *models.Mail) (*models.Mail, error)
	CreateMany(ctx context.Context, mails []models.Mail) (bool, error)
	// Update updates a mail object
	Update(ctx context.Context, mail *models.Mail) error
	// CleanupTrash cleans up trash mails
	CleanupTrash(ctx context.Context, userID *primitive.ObjectID, days *int) error
}

// MailRepository handles database operations related to mails
type MailRepository struct {
	collection *mongo.Collection
}

// NewMailRepository creates a new mail repository instance
func NewMailRepository(database *mongo.Database) MailRepositoryInterface {
	if database == nil {
		database = db.Database
	}
	return &MailRepository{
		collection: database.Collection(mailCollection),
	}
}

// GetAll retrieves mails for a user. If page and limit are >0, returns paginated results and total count. If page or limit <=0, returns all mails and total count.
func (r *MailRepository) GetAll(ctx context.Context, userID primitive.ObjectID, page, limit int64) ([]*models.Mail, int64, error) {
	filter := bson.M{"user_id": userID}
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build find options: always sort by created_at desc to return most recent first
	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	if page > 0 && limit > 0 {
		skip := (page - 1) * limit
		findOpts.SetSkip(skip)
		findOpts.SetLimit(limit)
	}

	cursor, err := r.collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var mails []*models.Mail
	if err = cursor.All(ctx, &mails); err != nil {
		return nil, 0, err
	}

	return mails, totalCount, nil
}

// GetByID retrieves a mail by its ID
func (r *MailRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Mail, error) {
	filter := bson.M{"_id": id}

	var mail models.Mail
	err := r.collection.FindOne(ctx, filter).Decode(&mail)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &mail, nil
}

// Create creates a new mail
func (r *MailRepository) Create(ctx context.Context, mail *models.Mail) (*models.Mail, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	if mail.ID == nil {
		id := primitive.NewObjectID()
		mail.ID = &id
	}

	// Only set CreatedAt if it hasn't been provided by the caller (tests may set it for deterministic ordering)
	if mail.CreatedAt == nil {
		mail.CreatedAt = &now
	}
	// Always update UpdatedAt to now
	mail.UpdatedAt = &now

	_, err := r.collection.InsertOne(ctx, mail)
	if err != nil {
		return nil, err
	}

	return mail, nil
}

// CreateMany creates multiple mails
func (r *MailRepository) CreateMany(ctx context.Context, mails []models.Mail) (bool, error) {
	now := primitive.NewDateTimeFromTime(time.Now())

	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return false, err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		for _, mail := range mails {
			id := primitive.NewObjectID()
			mail.ID = &id
			mail.CreatedAt = &now
			mail.UpdatedAt = &now

			_, err := r.collection.InsertOne(sessCtx, mail)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// Update updates a mail object
func (r *MailRepository) Update(ctx context.Context, mail *models.Mail) error {
	if mail == nil || mail.ID == nil {
		return nil
	}

	now := primitive.NewDateTimeFromTime(time.Now())
	mail.UpdatedAt = &now

	filter := bson.M{"_id": *mail.ID}
	update := bson.M{
		"$set": bson.M{
			"read":       mail.Read,
			"archived":   mail.Archived,
			"trashed":    mail.Trashed,
			"trashed_at": mail.TrashedAt,
			"updated_at": mail.UpdatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// CleanupTrash cleans up trash mails when trashed_at is 30 days or older and trashed is true
func (r *MailRepository) CleanupTrash(ctx context.Context, userID *primitive.ObjectID, days *int) error {
	daysAgo := time.Now().AddDate(0, 0, -30)

	// if days is -1, then delete all trashed mails now
	if days != nil && *days == -1 {
		daysAgo = time.Now()
	} else if days != nil && *days > 0 {
		daysAgo = time.Now().AddDate(0, 0, -*days)
	}
	filter := bson.M{
		"trashed":    true,
		"trashed_at": bson.M{"$lte": primitive.NewDateTimeFromTime(daysAgo)},
	}
	if userID != nil {
		filter["user_id"] = userID
	}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

