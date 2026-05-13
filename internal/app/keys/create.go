package keys

import (
	"encoding/base64"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) CreateKeys(user_id int, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error) {
	_, err := base64.StdEncoding.DecodeString(keys.Ciphertext)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.InvalidData("invalid ciphertext")
	}

	nonceBytes, err := base64.StdEncoding.DecodeString(keys.Nonce)
	if err != nil || len(nonceBytes) != 12 {
		if err != nil {
			logger.LogError(err.Error(), "keys-app")
		}
		return nil, domain.InvalidData("invalid nonce")
	}

	saltBytes, err := base64.StdEncoding.DecodeString(keys.Salt)
	if err != nil || len(saltBytes) != 16 {
		if err != nil {
			logger.LogError(err.Error(), "keys-app")
		}
		return nil, domain.InvalidData("invalid salt")
	}

	_, err = k.keys.GetByUserID(user_id, keys.Type)
	if err == nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("keys is already exists")
	}

	keys, err = k.keys.Create(&domain.EncryptedKeys{
		UserID:     user_id,
		Ciphertext: keys.Ciphertext,
		Nonce:      keys.Nonce,
		Salt:       keys.Salt,
	})
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("failed to save keys")
	}

	return keys, nil
}
