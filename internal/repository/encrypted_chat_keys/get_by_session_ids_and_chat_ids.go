package keys

import (
	"time"

	"github.com/lib/pq"
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) GetBySessionIDsAndChatIDs(sessionIDs []int, chatIDs []int) ([]*domain.EncryptedChatKeys, error) {
	if len(sessionIDs) == 0 || len(chatIDs) == 0 {
		return []*domain.EncryptedChatKeys{}, nil
	}

	var keys []*domain.EncryptedChatKeys

	query := `
		SELECT 
			id, chat_id, session_id, from_session_id, encrypted_key, encapsulated_key, 
			cek_wrap, cek_wrap_iv, salt, nonce, created_at
		FROM encrypted_chat_keys
		WHERE chat_id = ANY($1) AND session_id = ANY($2)
	`

	start := time.Now()
	err := r.db.Select(&keys, query, pq.Array(chatIDs), pq.Array(sessionIDs))
	duration := time.Since(start)
	metrics.ObserveDB("encrypted_chat_keys_get_by_session_and_chat_ids", duration, err)

	if err != nil {
		return nil, err
	}

	return keys, nil
}
