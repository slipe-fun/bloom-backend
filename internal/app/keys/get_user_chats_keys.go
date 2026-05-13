package keys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) GetUserKeys(user_id int, keys_type string) (*domain.EncryptedKeys, error) {
	keys, err := k.keys.GetByUserID(user_id, keys_type)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("failed to get keys")
	}

	return keys, nil
}
