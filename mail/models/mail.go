package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MailAttachment represents a file attachment
type MailAttachment struct {
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	StoragePath string `bson:"storage_path" json:"storage_path"` // S3 storage key/reference
	StorageType string `bson:"storage_type" json:"storage_type"` // S3, MongoDB, etc.
	Size        int64  `bson:"size" json:"size"`                 // File size in bytes
}

// Mail represents a mail message
type Mail struct {
	ID             primitive.ObjectID  `bson:"_id" json:"id"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Headers        interface{}         `bson:"headers" json:"headers"`
	TextContent    string              `bson:"text_content" json:"text_content"`
	HTMLContent    string              `bson:"html_content" json:"html_content"`
	Attachments    []MailAttachment    `bson:"attachments" json:"attachments"`
	Archived       bool                `bson:"archived" json:"archived"`
	Trashed        bool                `bson:"trashed" json:"trashed"`
	Greylisted     bool                `bson:"graylisted" json:"graylisted"`
	Rejected       bool                `bson:"rejected" json:"rejected"`
	RewriteSubject bool                `bson:"rewrite_subject" json:"rewrite_subject"`
	CreatedAt      *primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt      *primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
