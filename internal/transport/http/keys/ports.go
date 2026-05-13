package keys

import "github.com/slipe-fun/skid-backend/internal/domain"

type KeysApp interface {
	CreateKeys(user_id int, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
	GetUserKeys(user_id int, keys_type string) (*domain.EncryptedKeys, error)
}
