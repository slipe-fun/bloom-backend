package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/crypto"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) CreateChat(tokenStr string, recipient int) (*domain.Chat, *domain.Session, error) {
	session, err := c.sessionApp.GetSession(tokenStr)
	if err != nil {
		return nil, nil, err
	}

	encryptionKey, err := crypto.GenerateEncryptionKey()
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, nil, domain.Failed("failed to generate encryption key")
	}

	encKey := encryptionKey
	chat, err := c.chats.Create(&domain.Chat{
		Members: []domain.Member{
			{
				ID:             session.UserID,
				KyberPublicKey: "",
				EcdhPublicKey:  "",
				EdPublicKey:    "",
			},
			{
				ID:             recipient,
				KyberPublicKey: "",
				EcdhPublicKey:  "",
				EdPublicKey:    "",
			},
		},
		EncryptionKey: &encKey,
	})

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, nil, domain.Failed("failed to create chat")
	}

	return chat, session, nil
}
