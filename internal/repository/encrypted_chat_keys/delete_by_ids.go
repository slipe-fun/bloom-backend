package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) DeleteByIDs(ids []int) error {
	if len(ids) == 0 {
		return nil
	}

	query := `
		DELETE FROM encrypted_chat_keys
		WHERE id = ANY($1)
	`

	start := time.Now()
	_, err := r.db.Exec(query, ids)
	duration := time.Since(start)
	metrics.ObserveDB("encrypted_chat_keys_delete_by_ids", duration, err)

	return err
}
