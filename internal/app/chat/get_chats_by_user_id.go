package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatsByUserID(user_id int) ([]*domain.ChatWithLastMessage, error) {
	chats, err := c.chats.GetByUserID(user_id)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chats not found")
	}

	return chats, nil
}
