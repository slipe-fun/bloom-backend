package keys

import "github.com/slipe-fun/skid-backend/internal/domain"

type KeysApp interface {
	CreateKeys(user_id int, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
	GetUserChatKeys(user_id int) (*domain.EncryptedKeys, error)
}
