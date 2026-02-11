package chat

import (
	"github.com/slipe-fun/skid-backend/internal/domain"
	"github.com/slipe-fun/skid-backend/internal/pkg/logger"
)

func (c *ChatApp) GetChatsByUserID(token string) ([]*domain.ChatWithLastMessage, error) {
	session, err := c.sessionApp.GetSession(token)
	if err != nil {
		return nil, err
	}

	chats, err := c.chats.GetByUserID(session.UserID)

	if err != nil {
		logger.LogError(err.Error(), "chat-app")
		return nil, domain.NotFound("chats not found")
	}

	return chats, nil
}
