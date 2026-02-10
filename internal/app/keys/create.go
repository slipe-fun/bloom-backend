package keys

import (
	"encoding/base64"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) CreateKeys(tokenStr string, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error) {
	session, err := k.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	_, err = base64.StdEncoding.DecodeString(keys.Ciphertext)
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

	existingKeys, err := k.keys.GetByUserId(session.UserID)
	if err == nil {
		existingKeys.Ciphertext = keys.Ciphertext
		existingKeys.Nonce = keys.Nonce
		existingKeys.Salt = keys.Salt

		err = k.keys.Edit(existingKeys)
		if err != nil {
			logger.LogError(err.Error(), "keys-app")
			return nil, domain.Failed("failed to save keys")
		}
		return existingKeys, nil
	}

	keys, err = k.keys.Create(&domain.EncryptedKeys{
		UserID:     session.UserID,
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
