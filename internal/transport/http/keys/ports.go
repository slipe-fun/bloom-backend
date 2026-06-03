package keys

import "github.com/slipe-fun/skid-backend/internal/domain"

type KeysApp interface {
	GetUserKeys(user_id int, keys_type string) (*domain.EncryptedKeys, error)
}
