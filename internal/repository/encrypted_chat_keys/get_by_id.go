package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) GetByID(id int) (*domain.EncryptedChatKeys, error) {
	var keys domain.EncryptedChatKeys

	query := `SELECT id, chat_id, session_id, encrypted_key, encapsulated_key, nonce, salt, created_at FROM encrypted_chat_keys WHERE id = $1`

	start := time.Now()

	err := r.db.Get(&keys, query, id)

	duration := time.Since(start)

	metrics.ObserveDB("encrypted_chat_keys_get_by_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &keys, nil
}
