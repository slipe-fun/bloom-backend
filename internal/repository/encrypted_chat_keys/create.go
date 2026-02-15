package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) Create(keys *domain.EncryptedChatKeys) (*domain.EncryptedChatKeys, error) {
	query := `INSERT INTO encrypted_chat_keys (chat_id, session_id, encrypted_key, encapsulated_key, nonce, salt)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id, chat_id, session_id, encrypted_key, encapsulated_key, nonce, salt, created_at`

	var created domain.EncryptedChatKeys

	start := time.Now()

	err := r.db.QueryRow(query, keys.ChatID, keys.SessionID, keys.EncryptedKey, keys.EncapsulatedKey, keys.Nonce, keys.Salt).
		Scan(&created.ID, &created.ChatID, &created.SessionID, &created.EncryptedKey, &created.EncapsulatedKey, &created.Nonce, &created.Salt, &created.CreatedAt)

	duration := time.Since(start)

	metrics.ObserveDB("encrypted_chat_keys_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
