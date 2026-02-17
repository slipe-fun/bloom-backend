package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/domain"

type EncryptedChatKeysRepo interface {
	Create(keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, error)
}

type SessionApp interface {
	GetSessionByID(id int) (*domain.Session, error)
	GetSessionByIDs(ids []int) ([]*domain.Session, error)
}

type ChatRepo interface {
	GetByID(id int) (*domain.Chat, error)
}
