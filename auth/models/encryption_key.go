package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// EncryptionKey represents the structure for encryption keys for a user
type EncryptionKey struct {
	UserKey      string              `json:"userKey" bson:"user_key"`
	BackupKey    string              `json:"backupKey" bson:"backup_key"`
	Salt         string              `json:"salt" bson:"salt"`
	MnemonicSalt string              `json:"mnemonicSalt" bson:"mnemonic_salt"`
	Type         *string             `json:"type" bson:"type"`
	CreatedAt    *primitive.DateTime `json:"createdAt" bson:"created_at"`
	UpdatedAt    *primitive.DateTime `json:"updatedAt" bson:"updated_at"`
}
