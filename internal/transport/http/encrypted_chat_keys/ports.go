package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/domain"

type EncryptedChatKeysApp interface {
	AddKeys(userID, chatID int, keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, error)
}
