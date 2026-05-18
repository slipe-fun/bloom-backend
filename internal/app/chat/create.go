package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) CreateChat(user_id, recipient int, handshake domain.Handshake) (*domain.Chat, error) {
	chat, err := c.chats.Create(&domain.Chat{
		Members: []domain.Member{
			{
				ID: user_id,
			},
			{
				ID: recipient,
			},
		},
		Handshake: &handshake,
	})

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.Failed("failed to create chat")
	}

	return chat, nil
}
