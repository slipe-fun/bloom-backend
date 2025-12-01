package KeysRepo

import "github.com/slipe-fun/skid-backend/internal/domain"

func (k *KeysRepo) Create(keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error) {
	query := `INSERT INTO keys (user_id, ciphertext, nonce, salt) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id, user_id, ciphertext, nonce, salt`

	var created domain.EncryptedKeys
	err := k.db.QueryRow(query, keys.UserID, keys.Ciphertext, keys.Nonce, keys.Salt).
		Scan(&created.ID, &created.UserID, &created.Ciphertext, &created.Nonce, &created.Salt)

	if err != nil {
		return nil, err
	}

	return &created, nil
}
