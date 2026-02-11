package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatWithUsers(user_id int, recipient int) (*domain.Chat, error) {
	chat, err := c.chats.GetWithUsers(user_id, recipient)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotExpired("chats not found")
	}

	return chat, nil
}
