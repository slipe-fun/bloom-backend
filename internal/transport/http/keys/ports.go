package keys

import "github.com/slipe-fun/skid-backend/internal/domain"

type KeysApp interface {
	CreateKeys(token string, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error)
	GetUserChatsKeys(token string) (*domain.EncryptedKeys, error)
}
