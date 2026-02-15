package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (r *EncryptedChatKeysRepo) Delete(id int) error {
	query := `DELETE FROM encrypted_chat_keys WHERE id = $1`

	start := time.Now()

	_, err := r.db.Exec(query, id)

	duration := time.Since(start)

	metrics.ObserveDB("encrypted_chat_keys_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
