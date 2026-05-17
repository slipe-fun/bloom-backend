package keys

import (
	"encoding/base64"

	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) CreateKeys(user_id int, keys *domain.EncryptedKeys) (*domain.EncryptedKeys, error) {
	user, err := k.users.GetByID(user_id)
	if err != nil {
		return nil, domain.Failed("database error")
	}
	if user == nil {
		return nil, domain.Failed("user not found")
	}

	_, err = k.keys.GetByUserID(user_id, keys.Type)
	if err == nil {
		return nil, domain.Failed("keys is already exists")
	}

	if len(keys.Ciphertext) > 10000 {
		return nil, domain.InvalidData("ciphertext too large")
	}
	_, err = base64.StdEncoding.DecodeString(keys.Ciphertext)
	if err != nil {
		return nil, domain.InvalidData("invalid ciphertext format")
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

	signatureBytes, err := base64.StdEncoding.DecodeString(keys.Signature)
	if err != nil || len(signatureBytes) != 114 {
		return nil, domain.InvalidData("invalid signature length")
	}

	edPublicKey, err := base64.StdEncoding.DecodeString(user.EdPublicKey)
	if err != nil || len(edPublicKey) != 57 {
		return nil, domain.InvalidData("invalid user public key")
	}

	if isValid, err := crypto.VerifyEncryptedMasterKeySignature(edPublicKey, signatureBytes, keys.Ciphertext, keys.Nonce, keys.Salt); !isValid {
		return nil, domain.InvalidData("invalid signature")
	} else if err != nil {
		return nil, domain.Failed("failed to verify signature")
	}

	keys, err = k.keys.Create(&domain.EncryptedKeys{
		UserID:     user_id,
		Type:       keys.Type,
		Ciphertext: keys.Ciphertext,
		Nonce:      keys.Nonce,
		Salt:       keys.Salt,
		Signature:  keys.Signature,
	})
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("failed to save keys")
	}

	return keys, nil
}
