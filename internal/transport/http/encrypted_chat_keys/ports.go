package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/domain"

type EncryptedChatKeysApp interface {
	AddKeys(userID, chatID int, keys []*domain.EncryptedChatKeys) ([]*domain.EncryptedChatKeys, int, error)
	GetBySessionID(session_id int) ([]*domain.EncryptedChatKeys, error)
}
