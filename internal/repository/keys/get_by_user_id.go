package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (k *KeysRepo) GetByUserID(user_id int, keys_type string) (*domain.EncryptedKeys, error) {
	var keys domain.EncryptedKeys

	query := `SELECT user_id, type, ciphertext, nonce, salt FROM keys WHERE user_id = $1 AND type = $2`

	start := time.Now()

	err := k.db.Get(&keys, query, user_id, keys_type)

	duration := time.Since(start)

	metrics.ObserveDB("keys_get_by_user_id", duration, err)

	if err != nil {
		return nil, err
	}

	return &keys, nil
}
