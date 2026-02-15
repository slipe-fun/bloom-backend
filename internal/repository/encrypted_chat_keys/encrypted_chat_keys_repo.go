package keys

import "github.com/jmoiron/sqlx"

type EncryptedChatKeysRepo struct {
	db *sqlx.DB
}

func NewEncryptedChatKeysRepo(db *sqlx.DB) *EncryptedChatKeysRepo {
	return &EncryptedChatKeysRepo{
		db: db,
	}
}
