package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatByID(user_id int, id int) (*domain.Chat, error) {
	chat, err := c.chats.GetByID(id)
	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chat not found")
	}

	if !c.HasMember(chat, user_id) {
		return nil, domain.NotFound("chat not found")
	}

	return chat, nil
}
