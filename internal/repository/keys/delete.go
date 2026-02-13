package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (k *KeysRepo) Delete(user_id int) error {
	query := `DELETE FROM keys WHERE user_id = $1`

	start := time.Now()

	_, err := k.db.Exec(query, user_id)

	duration := time.Since(start)

	metrics.ObserveDB("keys_delete", duration, err)

	if err != nil {
		return err
	}

	return nil
}
