package models

type EncryptionKey struct {
	UserKey string `json:"userKey" bson:"user_key"`
	BackupKey string `json:"backupKey" bson:"backup_key"`
	UserSalt string `json:"userSalt" bson:"user_salt"`
}