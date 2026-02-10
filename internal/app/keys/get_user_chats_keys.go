package keys

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (k *KeysApp) GetUserChatsKeys(tokenStr string) (*domain.EncryptedKeys, error) {
	session, err := k.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	keys, err := k.keys.GetByUserId(session.UserID)
	if err != nil {
		logger.LogError(err.Error(), "keys-app")
		return nil, domain.Failed("failed to get keys")
	}

	return keys, nil
}
