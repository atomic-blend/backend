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

	var opts []*options.FindOptions
	if page > 0 && limit > 0 {
		skip := (page - 1) * limit
		opts = append(opts, &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		})
	}

	var cursor *mongo.Cursor
	if len(opts) > 0 {
		cursor, err = r.collection.Find(ctx, filter, opts[0])
	} else {
		cursor, err = r.collection.Find(ctx, filter)
	}
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

	mail.CreatedAt = &now
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
