package keys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) GetUserChatsKeys(token string) (*domain.EncryptedKeys, error) {
	session, err := k.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	keys, err := k.keys.GetByUserID(session.UserID)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("failed to get keys")
	}

	return keys, nil
}
