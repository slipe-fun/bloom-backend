package encryptedchatkeys

import "github.com/slipe-fun/skid-backend/internal/domain"

type EncryptedChatKeysApp interface {
	AddKeys(userID int, batch []domain.AddKeyBatchItem) ([]*domain.EncryptedChatKeys, error)
	GetBySessionID(session_id int) ([]*domain.EnrichedChatKey, error)
}
