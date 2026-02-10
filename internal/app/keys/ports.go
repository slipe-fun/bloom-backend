package keys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type SessionApp interface {
	GetSession(token string) (*domain.Session, error)
}

type KeysRepo interface {
	GetByUserId(user_id int) (*domain.EncryptedKeys, error)
	Edit(keys *domain.EncryptedKeys) error
	Create(keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
}
