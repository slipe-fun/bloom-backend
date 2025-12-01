package KeysRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (k *KeysRepo) GetByUserId(user_id int) (*domain.EncryptedKeys, error) {
	var keys domain.EncryptedKeys

	query := `SELECT user_id, ciphertext, nonce, salt FROM keys WHERE user_id = $1`
	err := k.db.Get(&keys, query, user_id)

	if err != nil {
		return nil, err
	}

	return &keys, nil
}
