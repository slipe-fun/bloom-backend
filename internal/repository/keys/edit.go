package keys

import "github.com/slipe-fun/skid-backend/internal/domain"

func (k *KeysRepo) Edit(keys *domain.EncryptedKeys) error {
	query := `
        UPDATE keys
        SET ciphertext = $1,
            nonce = $2,
			salt = $3
        WHERE user_id = $4
        RETURNING id, user_id, ciphertext, nonce, salt
    `

	_, err := k.db.Exec(query, keys.Ciphertext, keys.Nonce, keys.Salt, keys.UserID)

	if err != nil {
		return err
	}

	return nil
}
