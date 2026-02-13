package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (k *KeysRepo) Edit(keys *domain.EncryptedKeys) error {
	query := `
        UPDATE keys
        SET ciphertext = $1,
            nonce = $2,
			salt = $3
        WHERE user_id = $4
        RETURNING id, user_id, ciphertext, nonce, salt
    `

	start := time.Now()

	_, err := k.db.Exec(query, keys.Ciphertext, keys.Nonce, keys.Salt, keys.UserID)

	duration := time.Since(start)

	metrics.ObserveDB("keys_edit", duration, err)

	if err != nil {
		return err
	}

	return nil
}
