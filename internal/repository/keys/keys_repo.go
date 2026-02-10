package keys

import (
	"github.com/jmoiron/sqlx"
	ChatRepo "github.com/slipe-fun/skid-backend/internal/repository/chat"
)

type KeysRepo struct {
	db       *sqlx.DB
	chatRepo *ChatRepo.ChatRepo
}

func NewKeysRepo(db *sqlx.DB, chatRepo *ChatRepo.ChatRepo) *KeysRepo {
	return &KeysRepo{db: db, chatRepo: chatRepo}
}
