package KeysApp

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
)

func (k *KeysApp) GetUserChatsKeys(tokenStr string) (*domain.EncryptedKeys, error) {
	session, err := k.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, err
	}

	keys, err := k.keys.GetByUserId(session.UserID)
	if err != nil {
		return nil, domain.Failed("failed to get keys")
	}

	return keys, nil
}
