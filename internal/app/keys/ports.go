package keys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

type KeysRepo interface {
	GetByUserID(user_id int, keys_type string) (*domain.EncryptedKeys, error)
	Create(keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
}
