package keys

import (
	"time"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/metrics"
)

func (k *KeysRepo) Create(keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error) {
	query := `INSERT INTO keys (user_id, type, ciphertext, nonce, salt, signature) 
	          VALUES ($1, $2, $3, $4, $5, $6) 
	          RETURNING id, user_id, type, ciphertext, nonce, salt, signature`

	var created domain.EncryptedKeys

	start := time.Now()

	err := k.db.QueryRow(query, keys.UserID, keys.Type, keys.Ciphertext, keys.Nonce, keys.Salt, keys.Signature).
		Scan(&created.ID, &created.UserID, &created.Type, &created.Ciphertext, &created.Nonce, &created.Salt, &created.Signature)

	duration := time.Since(start)

	metrics.ObserveDB("keys_create", duration, err)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
