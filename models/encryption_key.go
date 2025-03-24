package models

type EncryptionKey struct {
	UserKey string `json:"userKey" bson:"user_key"`
	BackupKey string `json:"backupKey" bson:"backup_key"`
	Salt string `json:"salt" bson:"salt"`
	MnemonicSalt string `json:"mnemonicSalt" bson:"mnemonic_salt"`
}